package services

import (
	"accommodations-service/client"
	"accommodations-service/config"
	"accommodations-service/domain"
	"accommodations-service/errors"
	"accommodations-service/orchestrator"
	"accommodations-service/repository"
	"accommodations-service/utils"
	"context"
	events "example/saga/create_accommodation"
	"fmt"
	"log"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/trace"
)

type AccommodationService struct {
	accommodationRepository *repository.AccommodationRepo
	validator               *utils.Validator
	reservationsClient      *client.ReservationsClient
	userClient              *client.UserClient
	fileStorage             *repository.FileStorage
	cache                   *repository.ImageCache
	orchestrator            *orchestrator.CreateAccommodationOrchestrator
	tracer                  trace.Tracer
	logger                  *config.Logger
}

func NewAccommodationService(accommodationRepo *repository.AccommodationRepo, validator *utils.Validator, reservationsClient *client.ReservationsClient, userClient *client.UserClient, fileStorage *repository.FileStorage, cache *repository.ImageCache, orchestrator *orchestrator.CreateAccommodationOrchestrator, tracer trace.Tracer, logger *config.Logger) *AccommodationService {
	return &AccommodationService{
		accommodationRepository: accommodationRepo,
		validator:               validator,
		reservationsClient:      reservationsClient,
		userClient:              userClient,
		fileStorage:             fileStorage,
		cache:                   cache,
		orchestrator:            orchestrator,
		tracer:                  tracer,
		logger:                  logger,
	}
}

func (as *AccommodationService) CreateAccommodation(accommodation domain.CreateAccommodation, image multipart.File, ctx context.Context) (*domain.AccommodationDTO, *errors.ErrorStruct) {
	ctx, span := as.tracer.Start(ctx, "AccommodationService.CreateAccommodation")
	defer span.End()
	var imageIds []string
	accomm := domain.Accommodation{
		Name:             accommodation.Name,
		Address:          accommodation.Address,
		City:             accommodation.City,
		Country:          accommodation.Country,
		UserName:         accommodation.UserName,
		UserId:           accommodation.UserId,
		Email:            accommodation.Email,
		Conveniences:     accommodation.Conveniences,
		MinNumOfVisitors: accommodation.MinNumOfVisitors,
		MaxNumOfVisitors: accommodation.MaxNumOfVisitors,
		Paying:           accommodation.Paying,
	}
	as.validator.ValidateAccommodation(&accomm)
	//as.validator.ValidateAvailabilities(&accommodation)
	validatorErrors := as.validator.GetErrors()
	if len(validatorErrors) > 0 {
		var constructedError string
		for _, message := range validatorErrors {
			constructedError += message + "\n"
			as.logger.LogError("accommodation-service", fmt.Sprintf("Errors in validating accommodations:"+message))
		}
		as.validator.ClearErrors()
		as.logger.LogError("accommodation-service", fmt.Sprintf("Bad password for user %v", constructedError))
		return nil, errors.NewError(constructedError, 400)
	}

	log.Println(accomm)
	uuidStr := uuid.New().String()
	imageIds = append(imageIds, uuidStr)
	as.fileStorage.WriteFile(ctx, image, uuidStr)
	as.cache.Post(ctx, image, uuidStr)
	accomm.ImageIds = imageIds
	accomm.Status = "Pending"
	newAccommodation, foundErr := as.accommodationRepository.SaveAccommodation(ctx, accomm)
	if foundErr != nil {
		as.logger.LogError("accommodation-service", fmt.Sprintf("Error saving accommodation"))
		as.logger.LogError("accommodation-service", fmt.Sprintf("Error:"+foundErr.GetErrorMessage()))
		return nil, foundErr
	}
	as.logger.LogInfo("accommodation-service", "New accommodation created with id "+newAccommodation.Id.Hex())
	id := newAccommodation.Id.Hex()

	//err := as.reservationsClient.SendCreatedReservationsAvailabilities(ctx, id, accommodation)
	reqData := domain.SendCreateAccommodationAvailability{
		AccommodationID: id,
		Location:        accommodation.Location,
		DateRange:       accommodation.AvailableAccommodationDates,
	}
	var eventsDateRangeCasted []events.AvailableAccommodationDates
	for _, value := range reqData.DateRange {
		val := events.AvailableAccommodationDates{
			AccommodationId: value.AccommodationId,
			Location:        value.Location,
			DateRange:       value.DateRange,
			Price:           value.Price,
		}
		eventsDateRangeCasted = append(eventsDateRangeCasted, val)
	}

	reqDataCasted := events.SendCreateAccommodationAvailability{
		AccommodationID: reqData.AccommodationID,
		Location:        reqData.Location,
		DateRange:       eventsDateRangeCasted,
	}
	err := as.orchestrator.Start(&reqDataCasted)
	if err != nil {
		as.logger.LogError("accommodation-service", fmt.Sprintf("Error in starting orchestrator"))
		as.logger.LogError("accommodation-service", fmt.Sprintf("Error:"+err.Error()))
		as.DeleteAccommodation(ctx, id)
		as.logger.LogInfo("accommodation-service", "Accommodation with id"+newAccommodation.Id.Hex()+"deleted")

		return nil, errors.NewError("Service is not responding correctly", 500)
	}

	return &domain.AccommodationDTO{
		Id:               id,
		Name:             accommodation.Name,
		UserName:         accommodation.UserName,
		UserId:           accommodation.UserId,
		Email:            accommodation.Email,
		Address:          accommodation.Address,
		City:             accommodation.City,
		Country:          accommodation.Country,
		Conveniences:     accommodation.Conveniences,
		MinNumOfVisitors: accommodation.MinNumOfVisitors,
		MaxNumOfVisitors: accommodation.MaxNumOfVisitors,
		ImageIds:         imageIds,
		Status:           accommodation.Status,
		Paying:           accommodation.Paying,
	}, nil
}

func (as *AccommodationService) GetImage(ctx context.Context, id string) ([]byte, *errors.ErrorStruct) {
	ctx, span := as.tracer.Start(ctx, "AccommodationService.GetImage")
	defer span.End()
	file, err := as.fileStorage.ReadFile(ctx, id)
	if err != nil {
		as.logger.LogError("accommodations-service", fmt.Sprintf("Unable to get image from file storage"))
		as.logger.LogError("accommodation-service", fmt.Sprintf("Error:"+err.Error()))
		return nil, errors.NewError("image read error", 500)
	}
	as.cache.Create(ctx, file, id)
	as.logger.LogInfo("accommodation-service", "Cache file with id"+id+"created")
	return file, nil
}

func (as *AccommodationService) GetCache(ctx context.Context, key string) ([]byte, error) {
	ctx, span := as.tracer.Start(ctx, "AccommodationService.GetCache")
	defer span.End()
	data, err := as.cache.Get(ctx, key)
	as.logger.LogInfo("accommodation-service", "Cache data retrieved successfully")
	return data, err
}

func (as *AccommodationService) GetAllAccommodations(ctx context.Context) ([]*domain.AccommodationDTO, *errors.ErrorStruct) {
	ctx, span := as.tracer.Start(ctx, "AccommodationService.GetAllAccommodations")
	defer span.End()
	accommodations, err := as.accommodationRepository.GetAllAccommodations(ctx)

	if err != nil {
		as.logger.LogError("accommodations-service", fmt.Sprintf("Unable to get all accommodations"))
		as.logger.LogError("accommodation-service", fmt.Sprintf("Error:"+err.GetErrorMessage()))
		return nil, err
	}

	var domainAccommodations []*domain.AccommodationDTO
	for _, accommodation := range accommodations {
		id := accommodation.Id.Hex()
		imageIds := accommodation.ImageIds

		domainAccommodations = append(domainAccommodations, &domain.AccommodationDTO{
			Id:               id,
			Name:             accommodation.Name,
			UserName:         accommodation.UserName,
			UserId:           accommodation.UserId,
			Email:            accommodation.Email,
			Address:          accommodation.Address,
			City:             accommodation.City,
			Country:          accommodation.Country,
			Conveniences:     accommodation.Conveniences,
			MinNumOfVisitors: accommodation.MinNumOfVisitors,
			MaxNumOfVisitors: accommodation.MaxNumOfVisitors,
			ImageIds:         imageIds,
			Rating:           accommodation.Rating,
			Status:           accommodation.Status,
			Paying:           accommodation.Paying,
		})
	}
	as.logger.LogInfo("accommodation-service", "Successfully retrieved all available accommodations")
	return domainAccommodations, nil
}
func (as *AccommodationService) GetAccommodationById(ctx context.Context, accommodationId string) (*domain.Accommodation, *errors.ErrorStruct) {
	ctx, span := as.tracer.Start(ctx, "AccommodationService.GetAccommodationById")
	defer span.End()
	accomm, err := as.accommodationRepository.GetAccommodationById(ctx, accommodationId)

	if err != nil {
		as.logger.LogError("accommodations-service", fmt.Sprintf("Unable to get accommodation with id %s", accommodationId))
		as.logger.LogError("accommodation-service", fmt.Sprintf("Error:"+err.GetErrorMessage()))
		return nil, err
	}
	id, _ := accomm.Id.MarshalJSON()
	as.logger.LogInfo("accommodation-service", "Successfully retrieved accommodation with id"+accommodationId)
	return &domain.Accommodation{
		Id:               primitive.ObjectID(id),
		Name:             accomm.Name,
		UserName:         accomm.UserName,
		UserId:           accomm.UserId,
		Email:            accomm.Email,
		Address:          accomm.Address,
		City:             accomm.City,
		Country:          accomm.Country,
		Conveniences:     accomm.Conveniences,
		MinNumOfVisitors: accomm.MinNumOfVisitors,
		MaxNumOfVisitors: accomm.MaxNumOfVisitors,
		ImageIds:         accomm.ImageIds,
		Status:           accomm.Status,
		Paying:           accomm.Paying,
	}, nil

}

func (as *AccommodationService) FindAccommodationByIds(ctx context.Context, ids []string) ([]*domain.AccommodationDTO, *errors.ErrorStruct) {
	ctx, span := as.tracer.Start(ctx, "AccommodationService.FindAccommodationByIds")
	defer span.End()
	accomm, err := as.accommodationRepository.FindAccommodationByIds(ctx, ids)
	if err != nil {
		as.logger.LogError("accommodations-service", fmt.Sprintf("Unable to get accommodation with multiple ids"))
		as.logger.LogError("accommodation-service", fmt.Sprintf("Error:"+err.GetErrorMessage()))
		return nil, err
	}
	var domainAccommodations []*domain.AccommodationDTO
	for _, accommodation := range accomm {

		imageIds := accommodation.ImageIds
		id := accommodation.Id.Hex()
		domainAccommodations = append(domainAccommodations, &domain.AccommodationDTO{
			Id:               id,
			Name:             accommodation.Name,
			UserName:         accommodation.UserName,
			UserId:           accommodation.UserId,
			Email:            accommodation.Email,
			Address:          accommodation.Address,
			City:             accommodation.City,
			Country:          accommodation.Country,
			Conveniences:     accommodation.Conveniences,
			MinNumOfVisitors: accommodation.MinNumOfVisitors,
			MaxNumOfVisitors: accommodation.MaxNumOfVisitors,
			ImageIds:         imageIds,
			Rating:           accommodation.Rating,
			Status:           accommodation.Status,
			Paying:           accommodation.Paying,
		})
	}
	as.logger.LogInfo("accommodation-service", "Successfully retrieved accommodations with multiple ids")
	return domainAccommodations, nil

}

func (as *AccommodationService) UpdateAccommodation(ctx context.Context, updatedAccommodation domain.Accommodation) (*domain.Accommodation, *errors.ErrorStruct) {
	ctx, span := as.tracer.Start(ctx, "AccommodationService.UpdateAccommodation")
	defer span.End()
	as.validator.ValidateAccommodation(&updatedAccommodation)
	validatorErrors := as.validator.GetErrors()
	if len(validatorErrors) > 0 {
		var constructedError string
		for _, message := range validatorErrors {
			constructedError += message + "\n"
		}
		as.logger.LogError("accommodations-service", fmt.Sprintf("Unable to validate updated accommodation"))
		as.logger.LogError("accommodation-service", fmt.Sprintf("Error:"+constructedError))
		return nil, errors.NewError(constructedError, 400)
	}

	log.Println("Prije update")
	_, updateErr := as.accommodationRepository.UpdateAccommodationById(ctx, updatedAccommodation)
	if updateErr != nil {
		as.logger.LogError("accommodations-service", fmt.Sprintf("Unable updated accommodation"))
		as.logger.LogError("accommodation-service", fmt.Sprintf("Error:"+updateErr.GetErrorMessage()))
		return nil, errors.NewError("Unable to update", 500)
	}
	log.Println("Poslije update")
	as.logger.LogInfo("accommodation-service", "Successfully updated accommodation")
	return &domain.Accommodation{
		Id:               updatedAccommodation.Id,
		Name:             updatedAccommodation.Name,
		UserName:         updatedAccommodation.UserName,
		UserId:           updatedAccommodation.UserId,
		Email:            updatedAccommodation.Email,
		Address:          updatedAccommodation.Address,
		City:             updatedAccommodation.City,
		Country:          updatedAccommodation.Country,
		Conveniences:     updatedAccommodation.Conveniences,
		MinNumOfVisitors: updatedAccommodation.MinNumOfVisitors,
		MaxNumOfVisitors: updatedAccommodation.MaxNumOfVisitors,
		Status:           updatedAccommodation.Status,
		Paying:           updatedAccommodation.Paying,
	}, nil
}

func (as *AccommodationService) DeleteAccommodation(ctx context.Context, accommodationID string) (*domain.Accommodation, *errors.ErrorStruct) {
	ctx, span := as.tracer.Start(ctx, "AccommodationService.DeleteAccommodation")
	defer span.End()
	// Assuming validation checks are not necessary for deletion

	existingAccommodation, foundErr := as.accommodationRepository.GetAccommodationById(ctx, accommodationID)
	if foundErr != nil {
		as.logger.LogError("accommodations-service", fmt.Sprintf("Unable to get accommodation by id"+accommodationID))
		as.logger.LogError("accommodation-service", fmt.Sprintf("Error:"+foundErr.GetErrorMessage()))
		return nil, foundErr
	}

	deleteErr := as.accommodationRepository.DeleteAccommodationById(ctx, accommodationID)
	if deleteErr != nil {
		as.logger.LogError("accommodations-service", fmt.Sprintf("Unable to delete accommodation by id"+accommodationID))
		as.logger.LogError("accommodation-service", fmt.Sprintf("Error:"+deleteErr.GetErrorMessage()))
		return nil, deleteErr
	}
	as.logger.LogInfo("accommodation-service", "Successfully deleted accommodation with id"+accommodationID)
	return existingAccommodation, nil
}

func (as *AccommodationService) DeleteAccommodationsByUserId(ctx context.Context, userID string) *errors.ErrorStruct {
	ctx, span := as.tracer.Start(ctx, "AccommodationService.DeleteAccommodationsByUserId")
	defer span.End()

	deleteErr := as.accommodationRepository.DeleteAccommodationsByUserId(ctx, userID)
	if deleteErr != nil {
		as.logger.LogError("accommodations-service", fmt.Sprintf("Unable to delete accommodation by user id:"+userID))
		as.logger.LogError("accommodation-service", fmt.Sprintf("Error:"+deleteErr.GetErrorMessage()))
		return deleteErr
	}

	as.logger.LogInfo("accommodation-service", "Successfully deleted accommodation with user id"+userID)

	return nil
}
func (as *AccommodationService) PutAccommodationRating(ctx context.Context, accommodationID string, accommodation domain.Accommodation) *errors.ErrorStruct {
	ctx, span := as.tracer.Start(ctx, "AccommodationService.PutAccommodationRating")
	defer span.End()

	err := as.accommodationRepository.PutAccommodationRating(ctx, accommodationID, accommodation.Rating)
	if err != nil {
		as.logger.LogError("accommodations-service", fmt.Sprintf("Unable to put accommodation rating into accommodation with id: "+accommodationID))
		as.logger.LogError("accommodation-service", fmt.Sprintf("Error:"+err.GetErrorMessage()))
		return errors.NewError("Error calling repository service", 500)
	}
	as.logger.LogInfo("accommodation-service", "Successfully added rating into accommodation with id:"+accommodationID)
	return nil
}

func (as *AccommodationService) SearchAccommodations(city, country string, numOfVisitors int, startDate string, endDate string, maxPrice int, conveniences []string, isDistinguishedString string, ctx context.Context) ([]domain.Accommodation, *errors.ErrorStruct) {
	ctx, span := as.tracer.Start(ctx, "AccommodationService.SearchAccommodations")
	defer span.End()
	log.Println("USLO U SERVIS")
	log.Println("Is distiguished na pocetku", isDistinguishedString)
	isDistinguished := false
	if isDistinguishedString == "true" {
		isDistinguished = true
	}

	log.Println("Start date", startDate)
	log.Println("EndDate", endDate)
	log.Println("Max Price", maxPrice)
	log.Println("isDistinguished", isDistinguished)

	accommodations, err := as.accommodationRepository.SearchAccommodations(ctx, city, country, numOfVisitors, maxPrice, conveniences)
	if err != nil {
		as.logger.LogError("accommodations-service", fmt.Sprintf("Unable to search accommodations"))
		as.logger.LogError("accommodation-service", fmt.Sprintf("Error:"+err.GetErrorMessage()))
		return nil, errors.NewError("Failed to find accommodations", 500) // Modify according to your error handling approach
	}
	var accommodationIDs []string

	for _, acc := range accommodations {
		accommodationIDs = append(accommodationIDs, acc.Id.Hex())

	}
	log.Println(accommodationIDs)
	if startDate == "" && endDate == "" && isDistinguished == false && maxPrice == 0 {
		return accommodations, nil
	}

	if startDate != "" && endDate != "" && isDistinguished == false && maxPrice == 0 {

		dateRange, err := generateDateRange(startDate, endDate)
		if err != nil {
			as.logger.LogError("accommodations-service", fmt.Sprintf("Unable to generate daterange"))
			as.logger.LogError("accommodation-service", fmt.Sprintf("Error:"+err.GetErrorMessage()))
			return nil, errors.NewError("Failed to generate dateRange", 500) // Modify according to your error handling approach
		}
		log.Println("dateRange je", dateRange)

		reservedIDs, err := as.reservationsClient.CheckAvailabilityForAccommodations(ctx, accommodationIDs, dateRange)
		if err != nil {
			return nil, errors.NewError("Failed to get reserved ids ", 500)
		}
		log.Println("Reservisani idevi", reservedIDs)
		log.Println("Sve nadjene akomodacije", accommodations)
		filteredAccommodations := removeAccommodations(accommodations, reservedIDs)
		as.logger.LogInfo("accommodation-service", "Successfully filtered accommodations")
		return filteredAccommodations, nil
	}

	if startDate != "" && endDate != "" && isDistinguished == true && maxPrice == 0 {

		dateRange, err := generateDateRange(startDate, endDate)
		if err != nil {
			as.logger.LogError("accommodations-service", fmt.Sprintf("Unable to generate daterange"))
			as.logger.LogError("accommodation-service", fmt.Sprintf("Error:"+err.GetErrorMessage()))
			return nil, errors.NewError("Failed to generate dateRange", 500) // Modify according to your error handling approach
		}

		reservedIDs, err := as.reservationsClient.CheckAvailabilityForAccommodations(ctx, accommodationIDs, dateRange)
		if err != nil {
			as.logger.LogError("accommodations-service", fmt.Sprintf("Unable check availability for accommodations"))
			as.logger.LogError("accommodation-service", fmt.Sprintf("Error:"+err.GetErrorMessage()))
			return nil, errors.NewError("Failed to get reserved ids ", 500)
		}
		log.Println("Reservisani idevi", reservedIDs)
		log.Println("Sve nadjene akomodacije", accommodations)
		filteredAccommodations := removeAccommodations(accommodations, reservedIDs)

		var distFiltered []domain.Accommodation

		for _, acc := range filteredAccommodations {
			user, _ := as.userClient.GetUserById(ctx, acc.UserId)
			log.Println("User je", user)

			if user.Distinguished == true {

				distFiltered = append(distFiltered, acc)
			}
		}
		as.logger.LogInfo("accommodation-service", "Successfully filtered accommodations")
		log.Println("Akomodacije sa distinguished likovima su", distFiltered)
		return distFiltered, nil

	}

	if startDate == "" && endDate == "" && isDistinguished == true && maxPrice == 0 {
		log.Println("USLO JE GDJE TREBA")

		var distFiltered []domain.Accommodation
		log.Println("Sve akomodacije", accommodations)

		for _, acc := range accommodations {
			log.Println("UserId:", acc.UserId)
			user, err := as.userClient.GetUserById(ctx, acc.UserId)
			log.Println("User je", user)
			if err != nil {
				as.logger.LogError("accommodations-service", fmt.Sprintf("Error getting user with id"+acc.UserId))
				as.logger.LogError("accommodation-service", fmt.Sprintf("Error:"+err.GetErrorMessage()))
				log.Println("Error getting user:", err)
				// Handle the error if needed
				continue
			}

			if user != nil && user.Distinguished == true {
				distFiltered = append(distFiltered, acc)
			}
		}
		as.logger.LogInfo("accommodation-service", "Successfully filtered accommodations")

		log.Println("Akomodacije sa distinguished likovima su", distFiltered)
		return distFiltered, nil
	}

	if startDate == "" && endDate == "" && isDistinguished == false && maxPrice != 0 {
		accBelowPrice, err := as.reservationsClient.GetAccommodationsBelowPrice(ctx, maxPrice)
		if err != nil {
			as.logger.LogError("accommodations-service", fmt.Sprintf("Error getting accommodation below the price of %d", maxPrice))
			as.logger.LogError("accommodation-service", fmt.Sprintf("Error:"+err.GetErrorMessage()))
			return nil, errors.NewError("Failed to get accommodations from reservations service", 500) // Modify according to your error handling approach
		}
		filteredAccommodationsWPrice := FilterAccommodationsByID(accBelowPrice, accommodations)
		as.logger.LogInfo("accommodation-service", "Successfully filtered accommodations")
		return filteredAccommodationsWPrice, nil
	}

	if startDate != "" && endDate != "" && isDistinguished == true && maxPrice != 0 {

		dateRange, err := generateDateRange(startDate, endDate)
		if err != nil {
			as.logger.LogError("accommodations-service", fmt.Sprintf("Error generating date range"))
			as.logger.LogError("accommodation-service", fmt.Sprintf("Error:"+err.GetErrorMessage()))
			return nil, errors.NewError("Failed to generate dateRange", 500) // Modify according to your error handling approach
		}
		log.Println("dateRange je", dateRange)

		reservedIDs, err := as.reservationsClient.CheckAvailabilityForAccommodations(ctx, accommodationIDs, dateRange)
		if err != nil {
			as.logger.LogError("accommodations-service", fmt.Sprintf("Error checking availabilities"))
			as.logger.LogError("accommodation-service", fmt.Sprintf("Error:"+err.GetErrorMessage()))
			return nil, errors.NewError("Failed to get reserved ids ", 500)
		}
		log.Println("Reservisani idevi", reservedIDs)
		log.Println("Sve nadjene akomodacije", accommodations)
		filteredAccommodations := removeAccommodations(accommodations, reservedIDs)

		var distFiltered []domain.Accommodation

		for _, acc := range filteredAccommodations {
			log.Println("UserId:", acc.UserId)
			user, err := as.userClient.GetUserById(ctx, acc.UserId)
			log.Println("User je", user)
			if err != nil {
				as.logger.LogError("accommodations-service", fmt.Sprintf("Error getting users for searching destinguished ones"))
				as.logger.LogError("accommodation-service", fmt.Sprintf("Error:"+err.GetErrorMessage()))
				log.Println("Error getting user:", err)
				// Handle the error if needed
				continue
			}

			if user != nil && user.Distinguished == true {
				distFiltered = append(distFiltered, acc)
			}
		}

		accBelowPrice, err := as.reservationsClient.GetAccommodationsBelowPrice(ctx, maxPrice)

		if err != nil {
			as.logger.LogError("accommodations-service", fmt.Sprintf("Error getting accommodation below price of %d", maxPrice))
			as.logger.LogError("accommodation-service", fmt.Sprintf("Error:"+err.GetErrorMessage()))
			return nil, errors.NewError("Failed to get accommodations from reservations service", 500) // Modify according to your error handling approach
		}
		filteredAccommodationsWPrice := FilterAccommodationsByID(accBelowPrice, distFiltered)
		as.logger.LogInfo("accommodation-service", "Successfully filtered accommodations")
		return filteredAccommodationsWPrice, nil

	}

	if startDate == "" && endDate == "" && isDistinguished == true && maxPrice != 0 {
		log.Println("UKUCAO SI DIST TRUE, MAX PRICE NIJE 0")

		var distFiltered []domain.Accommodation
		log.Println("Sve akomodacije", accommodations)

		for _, acc := range accommodations {
			log.Println("UserId:", acc.UserId)
			user, err := as.userClient.GetUserById(ctx, acc.UserId)
			log.Println("User je", user)
			if err != nil {
				as.logger.LogError("accommodations-service", fmt.Sprintf("Error getting user by id of %s", acc.UserId))
				as.logger.LogError("accommodation-service", fmt.Sprintf("Error:"+err.GetErrorMessage()))
				log.Println("Error getting user:", err)
				// Handle the error if needed
				continue
			}

			if user != nil && user.Distinguished == true {
				distFiltered = append(distFiltered, acc)
			}
		}

		accBelowPrice, err := as.reservationsClient.GetAccommodationsBelowPrice(ctx, maxPrice)
		if err != nil {
			as.logger.LogError("accommodations-service", fmt.Sprintf("Error getting accommodation below price of %d", maxPrice))
			as.logger.LogError("accommodation-service", fmt.Sprintf("Error:"+err.GetErrorMessage()))
			return nil, errors.NewError("Failed to get accommodations from reservations service", 500) // Modify according to your error handling approach
		}
		filteredAccommodationsWPrice := FilterAccommodationsByID(accBelowPrice, distFiltered)
		as.logger.LogInfo("accommodation-service", "Successfully filtered accommodations")
		return filteredAccommodationsWPrice, nil

	}

	return nil, errors.NewError("Failed to return anything", 500)

}

func FilterAccommodationsByID(ids []string, accommodations []domain.Accommodation) []domain.Accommodation {
	filteredAccommodations := make([]domain.Accommodation, 0)

	idSet := make(map[string]bool)
	for _, id := range ids {
		idSet[id] = true
	}

	for _, acc := range accommodations {
		if idSet[acc.Id.Hex()] {
			filteredAccommodations = append(filteredAccommodations, acc)
		}
	}

	return filteredAccommodations
}

func removeAccommodations(accommodations []domain.Accommodation, accommodationIDs []string) []domain.Accommodation {

	var filteredAccommodations []domain.Accommodation

	// Create a map for faster lookup of accommodationIDs
	idMap := make(map[string]bool)
	for _, id := range accommodationIDs {
		idMap[id] = true
	}

	// Check accommodations against accommodationIDs and remove if necessary
	for _, acc := range accommodations {
		if idMap[acc.Id.Hex()] {
			// If the ID exists in accommodationIDs, exclude it from filteredAccommodations
			continue
		}
		filteredAccommodations = append(filteredAccommodations, acc)
	}

	return filteredAccommodations
}

func generateDateRange(startDateStr, endDateStr string) ([]string, *errors.ErrorStruct) {
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		if err != nil {
			// Handle the error returned by the repository
			return nil, errors.NewError("Failed to parse date", 500) // Modify according to your error handling approach
		}
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		if err != nil {
			// Handle the error returned by the repository
			return nil, errors.NewError("Failed to parse date", 500) // Modify according to your error handling approach
		}
	}

	var dates []string
	for currentDate := startDate; !currentDate.After(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		dates = append(dates, currentDate.Format("2006-01-02"))
	}

	return dates, nil
}

func (as AccommodationService) ApproveAccommodation(ctx context.Context, id string) *errors.ErrorStruct {
	ctx, span := as.tracer.Start(ctx, "AccommodationService.ApproveAccommodation")
	defer span.End()
	log.Println("USLO DA POTVRDI AKOMODACIJU")
	acomm, err := as.GetAccommodationById(ctx, id)
	if err != nil {
		as.logger.LogError("accommodations-service", fmt.Sprintf("Failed to get accommodation by id in ApproveAccommodation func with id %s", id))
		as.logger.LogError("accommodation-service", fmt.Sprintf("Error:"+err.GetErrorMessage()))
		return err
	}
	acomm.Status = "Approved"
	err = as.accommodationRepository.PutAccommodationStatus(id, "Approved")
	if err != nil {
		as.logger.LogError("accommodations-service", fmt.Sprintf("Error setting accommodation status of %s", acomm.Status))
		as.logger.LogError("accommodation-service", fmt.Sprintf("Error:"+err.GetErrorMessage()))
		log.Println(err)
		return err
	}
	as.logger.LogInfo("accommodation-service", "Successfully put accommodation status")
	return nil
}

func (as AccommodationService) DenyAccommodation(ctx context.Context, id string) error {
	ctx, span := as.tracer.Start(ctx, "AccommodationService.DenyAccommodation")
	defer span.End()
	log.Println("DENY ACCOMMODATION")
	as.DeleteAccommodation(ctx, id)
	as.logger.LogInfo("accommodation-service", fmt.Sprintf("Accommodation with id %s deleted", id))
	return nil
}
