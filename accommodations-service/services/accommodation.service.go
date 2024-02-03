package services

import (
	"accommodations-service/client"
	"accommodations-service/domain"
	"accommodations-service/errors"
	"accommodations-service/repository"
	"accommodations-service/utils"
	"context"
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
	fileStorage             *repository.FileStorage
	cache                   *repository.ImageCache
	tracer                  trace.Tracer
}

func NewAccommodationService(accommodationRepo *repository.AccommodationRepo, validator *utils.Validator, reservationsClient *client.ReservationsClient, fileStorage *repository.FileStorage, cache *repository.ImageCache, tracer trace.Tracer) *AccommodationService {
	return &AccommodationService{
		accommodationRepository: accommodationRepo,
		validator:               validator,
		reservationsClient:      reservationsClient,
		fileStorage:             fileStorage,
		cache:                   cache,
		tracer:                  tracer,
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
	}
	as.validator.ValidateAccommodation(&accomm)
	//as.validator.ValidateAvailabilities(&accommodation)
	validatorErrors := as.validator.GetErrors()
	if len(validatorErrors) > 0 {
		var constructedError string
		for _, message := range validatorErrors {
			constructedError += message + "\n"
		}
		as.validator.ClearErrors()
		return nil, errors.NewError(constructedError, 400)
	}

	log.Println(accomm)
	uuidStr := uuid.New().String()
	imageIds = append(imageIds, uuidStr)
	as.fileStorage.WriteFile(ctx, image, uuidStr)
	as.cache.Post(ctx, image, uuidStr)
	accomm.ImageIds = imageIds
	newAccommodation, foundErr := as.accommodationRepository.SaveAccommodation(ctx, accomm)
	if foundErr != nil {
		return nil, foundErr
	}
	id := newAccommodation.Id.Hex()

	err := as.reservationsClient.SendCreatedReservationsAvailabilities(ctx, id, accommodation)
	if err != nil {
		as.DeleteAccommodation(ctx, id)
		return nil, errors.NewError("Service is not responding correcrtly", 500)
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
	}, nil
}

func (as *AccommodationService) GetImage(ctx context.Context, id string) ([]byte, *errors.ErrorStruct) {
	ctx, span := as.tracer.Start(ctx, "AccommodationService.GetImage")
	defer span.End()
	file, err := as.fileStorage.ReadFile(ctx, id)
	if err != nil {
		return nil, errors.NewError("image read error", 500)
	}
	as.cache.Create(ctx, file, id)
	return file, nil
}

func (as *AccommodationService) GetCache(ctx context.Context, key string) ([]byte, error) {
	ctx, span := as.tracer.Start(ctx, "AccommodationService.GetCache")
	defer span.End()
	data, err := as.cache.Get(ctx, key)
	return data, err
}

func (as *AccommodationService) GetAllAccommodations(ctx context.Context) ([]*domain.AccommodationDTO, *errors.ErrorStruct) {
	ctx, span := as.tracer.Start(ctx, "AccommodationService.GetAllAccommodations")
	defer span.End()
	accommodations, err := as.accommodationRepository.GetAllAccommodations(ctx)
	if err != nil {
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
		})
	}

	return domainAccommodations, nil
}
func (as *AccommodationService) GetAccommodationById(ctx context.Context, accommodationId string) (*domain.Accommodation, *errors.ErrorStruct) {
	ctx, span := as.tracer.Start(ctx, "AccommodationService.GetAccommodationById")
	defer span.End()
	accomm, err := as.accommodationRepository.GetAccommodationById(ctx, accommodationId)
	if err != nil {
		return nil, err
	}
	id, _ := accomm.Id.MarshalJSON()
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
	}, nil

}

func (as *AccommodationService) FindAccommodationByIds(ctx context.Context, ids []string) ([]*domain.AccommodationDTO, *errors.ErrorStruct) {
	ctx, span := as.tracer.Start(ctx, "AccommodationService.FindAccommodationByIds")
	defer span.End()
	accomm, err := as.accommodationRepository.FindAccommodationByIds(ctx, ids)
	if err != nil {
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
		})
	}
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
		return nil, errors.NewError(constructedError, 400)
	}

	log.Println("Prije update")
	_, updateErr := as.accommodationRepository.UpdateAccommodationById(ctx, updatedAccommodation)
	if updateErr != nil {
		return nil, errors.NewError("Unable to update", 500)
	}
	log.Println("Poslije update")

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
	}, nil
}

func (as *AccommodationService) DeleteAccommodation(ctx context.Context, accommodationID string) (*domain.Accommodation, *errors.ErrorStruct) {
	ctx, span := as.tracer.Start(ctx, "AccommodationService.DeleteAccommodation")
	defer span.End()
	// Assuming validation checks are not necessary for deletion

	existingAccommodation, foundErr := as.accommodationRepository.GetAccommodationById(ctx, accommodationID)
	if foundErr != nil {
		return nil, foundErr
	}

	deleteErr := as.accommodationRepository.DeleteAccommodationById(ctx, accommodationID)
	if deleteErr != nil {
		return nil, deleteErr
	}

	return existingAccommodation, nil
}

func (as *AccommodationService) DeleteAccommodationsByUserId(ctx context.Context, userID string) *errors.ErrorStruct {
	ctx, span := as.tracer.Start(ctx, "AccommodationService.DeleteAccommodationsByUserId")
	defer span.End()

	deleteErr := as.accommodationRepository.DeleteAccommodationsByUserId(ctx, userID)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}
func (as *AccommodationService) PutAccommodationRating(ctx context.Context, accommodationID string, accommodation domain.Accommodation) *errors.ErrorStruct {
	ctx, span := as.tracer.Start(ctx, "AccommodationService.PutAccommodationRating")
	defer span.End()

	err := as.accommodationRepository.PutAccommodationRating(ctx, accommodationID, accommodation.Rating)
	if err != nil {
		return errors.NewError("Error calling repository service", 500)
	}
	return nil
}

func (as *AccommodationService) SearchAccommodations(city, country string, numOfVisitors int, startDate string, endDate string, minPrice int, maxPrice int, conveniences []string, isDistinguished string, ctx context.Context) ([]domain.Accommodation, *errors.ErrorStruct) {
	ctx, span := as.tracer.Start(ctx, "AccommodationService.SearchAccommodations")
	defer span.End()
	log.Println("USLO U SERVIS")
	accommodations, err := as.accommodationRepository.SearchAccommodations(ctx, city, country, numOfVisitors, minPrice, maxPrice, conveniences)
	if err != nil {
		// Handle the error returned by the repository
		return nil, errors.NewError("Failed to find accommodations", 500) // Modify according to your error handling approach
	}
	var accommodationIDs []string
	for _, acc := range accommodations {
		accommodationIDs = append(accommodationIDs, acc.Id.Hex())
	}
	log.Println(accommodationIDs)
	if startDate == "" || endDate == "" {
		return accommodations, nil
	}
	dateRange, err := generateDateRange(startDate, endDate)
	if err != nil {
		// Handle the error returned by the repository
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
	log.Println("filtrirane akomodacije", filteredAccommodations)

	return filteredAccommodations, nil
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
