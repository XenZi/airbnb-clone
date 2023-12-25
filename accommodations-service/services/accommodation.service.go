package services

import (
	"accommodations-service/client"
	"accommodations-service/domain"
	"accommodations-service/errors"
	"accommodations-service/repository"
	"accommodations-service/utils"
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

type AccommodationService struct {
	accommodationRepository *repository.AccommodationRepo
	validator               *utils.Validator
	reservationsClient      *client.ReservationsClient
	fileStorage             *repository.FileStorage
}

func NewAccommodationService(accommodationRepo *repository.AccommodationRepo, validator *utils.Validator, reservationsClient *client.ReservationsClient, fileStorage *repository.FileStorage) *AccommodationService {
	return &AccommodationService{
		accommodationRepository: accommodationRepo,
		validator:               validator,
		reservationsClient:      reservationsClient,
		fileStorage:             fileStorage,
	}
}

func (as *AccommodationService) CreateAccommodation(accommodation domain.CreateAccommodation, ctx context.Context) (*domain.AccommodationDTO, *errors.ErrorStruct) {
	images := accommodation.Images
	var imageIds []string
	for _, val := range images {
		id := uuid.New().String()
		err := as.fileStorage.WriteFile(val, id)
		if err != nil {
			return nil, errors.NewError("internal image storage error", 500)
		}
		imageIds = append(imageIds, id)
	}

	accomm := domain.Accommodation{
		Name:             accommodation.Name,
		Address:          accommodation.Address,
		City:             accommodation.City,
		Country:          accommodation.Country,
		UserName:         accommodation.UserName,
		UserId:           accommodation.UserId,
		Conveniences:     accommodation.Conveniences,
		MinNumOfVisitors: accommodation.MinNumOfVisitors,
		MaxNumOfVisitors: accommodation.MaxNumOfVisitors,
		ImageIds:         imageIds,
	}
	as.validator.ValidateAccommodation(&accomm)
	as.validator.ValidateAvailabilities(&accommodation)
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
	newAccommodation, foundErr := as.accommodationRepository.SaveAccommodation(accomm)
	if foundErr != nil {
		return nil, foundErr
	}
	id := newAccommodation.Id.Hex()

	err := as.reservationsClient.SendCreatedReservationsAvailabilities(ctx, id, accommodation)
	if err != nil {
		as.DeleteAccommodation(id)
		return nil, errors.NewError("Service is not responding correcrtly", 500)
	}
	var imageBytes [][]byte
	for _, val := range imageIds {
		img, err := as.fileStorage.ReadFile(val)
		if err != nil {
			return nil, errors.NewError("image read error", 500)
		}
		imageBytes = append(imageBytes, img)
	}

	return &domain.AccommodationDTO{
		Id:               id,
		Name:             accommodation.Name,
		UserName:         accommodation.UserName,
		UserId:           accommodation.UserId,
		Address:          accommodation.Address,
		City:             accommodation.City,
		Country:          accommodation.Country,
		Conveniences:     accommodation.Conveniences,
		MinNumOfVisitors: accommodation.MinNumOfVisitors,
		MaxNumOfVisitors: accommodation.MaxNumOfVisitors,
		Images:           imageBytes,
	}, nil

}

func (as *AccommodationService) GetAllAccommodations() ([]*domain.AccommodationDTO, *errors.ErrorStruct) {
	accommodations, err := as.accommodationRepository.GetAllAccommodations()
	if err != nil {
		return nil, err
	}

	var domainAccommodations []*domain.AccommodationDTO
	for _, accommodation := range accommodations {
		id := accommodation.Id.Hex()
		imageIds := accommodation.ImageIds
		var images [][]byte
		for _, val := range imageIds {
			img, err := as.fileStorage.ReadFile(val)
			if err != nil {
				return nil, errors.NewError("image read error", 500)
			}
			images = append(images, img)
		}
		domainAccommodations = append(domainAccommodations, &domain.AccommodationDTO{
			Id:               id,
			Name:             accommodation.Name,
			UserName:         accommodation.UserName,
			UserId:           accommodation.UserId,
			Address:          accommodation.Address,
			City:             accommodation.City,
			Country:          accommodation.Country,
			Conveniences:     accommodation.Conveniences,
			MinNumOfVisitors: accommodation.MinNumOfVisitors,
			MaxNumOfVisitors: accommodation.MaxNumOfVisitors,
			Images:           images,
		})
	}

	return domainAccommodations, nil
}
func (as *AccommodationService) GetAccommodationById(accommodationId string) (*domain.Accommodation, *errors.ErrorStruct) {
	accomm, err := as.accommodationRepository.GetAccommodationById(accommodationId)
	if err != nil {
		return nil, err
	}
	id, _ := accomm.Id.MarshalJSON()
	return &domain.Accommodation{
		Id:               primitive.ObjectID(id),
		Name:             accomm.Name,
		UserName:         accomm.UserName,
		UserId:           accomm.UserId,
		Address:          accomm.Address,
		City:             accomm.City,
		Country:          accomm.Country,
		Conveniences:     accomm.Conveniences,
		MinNumOfVisitors: accomm.MinNumOfVisitors,
		MaxNumOfVisitors: accomm.MaxNumOfVisitors,
	}, nil

}

func (as *AccommodationService) UpdateAccommodation(updatedAccommodation domain.Accommodation) (*domain.Accommodation, *errors.ErrorStruct) {
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
	_, updateErr := as.accommodationRepository.UpdateAccommodationById(updatedAccommodation)
	if updateErr != nil {
		return nil, errors.NewError("Unable to update", 500)
	}
	log.Println("Poslije update")

	return &domain.Accommodation{
		Id:               updatedAccommodation.Id,
		Name:             updatedAccommodation.Name,
		UserName:         updatedAccommodation.UserName,
		UserId:           updatedAccommodation.UserId,
		Address:          updatedAccommodation.Address,
		City:             updatedAccommodation.City,
		Country:          updatedAccommodation.Country,
		Conveniences:     updatedAccommodation.Conveniences,
		MinNumOfVisitors: updatedAccommodation.MinNumOfVisitors,
		MaxNumOfVisitors: updatedAccommodation.MaxNumOfVisitors,
	}, nil
}

func (as *AccommodationService) DeleteAccommodation(accommodationID string) (*domain.Accommodation, *errors.ErrorStruct) {
	// Assuming validation checks are not necessary for deletion

	existingAccommodation, foundErr := as.accommodationRepository.GetAccommodationById(accommodationID)
	if foundErr != nil {
		return nil, foundErr
	}

	deleteErr := as.accommodationRepository.DeleteAccommodationById(accommodationID)
	if deleteErr != nil {
		return nil, deleteErr
	}

	return existingAccommodation, nil
}

func (as *AccommodationService) DeleteAccommodationsByUserId(userID string) *errors.ErrorStruct {

	deleteErr := as.accommodationRepository.DeleteAccommodationsByUserId(userID)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

func (as *AccommodationService) SearchAccommodations(city, country string, numOfVisitors int, startDate string, endDate string, ctx context.Context) ([]domain.Accommodation, *errors.ErrorStruct) {
	log.Println("USLO U SERVIS")
	accommodations, err := as.accommodationRepository.SearchAccommodations(city, country, numOfVisitors)
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
