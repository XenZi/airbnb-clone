package services

import (
	"accommodations-service/domain"
	"accommodations-service/errors"
	"accommodations-service/repository"
	"accommodations-service/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"strconv"
)

type AccommodationService struct {
	accommodationRepository *repository.AccommodationRepo
	validator               *utils.Validator
}

func NewAccommodationService(accommodationRepo *repository.AccommodationRepo, validator *utils.Validator) *AccommodationService {
	return &AccommodationService{
		accommodationRepository: accommodationRepo,
		validator:               validator,
	}
}

func (as *AccommodationService) CreateAccommodation(accommodation domain.Accommodation) (*domain.AccommodationDTO, *errors.ErrorStruct) {
	as.validator.ValidateAccommodation(&accommodation)
	validatorErrors := as.validator.GetErrors()
	if len(validatorErrors) > 0 {
		var constructedError string
		for _, message := range validatorErrors {
			constructedError += message + "\n"
		}
		return nil, errors.NewError(constructedError, 400)
	}
	accomm := domain.Accommodation{
		Name:             accommodation.Name,
		Location:         accommodation.Location,
		Conveniences:     accommodation.Conveniences,
		MinNumOfVisitors: accommodation.MinNumOfVisitors,
		MaxNumOfVisitors: accommodation.MaxNumOfVisitors,
	}
	newAccommodation, foundErr := as.accommodationRepository.SaveAccommodation(accomm)
	if foundErr != nil {
		return nil, foundErr
	}
	id := newAccommodation.Id.Hex()

	return &domain.AccommodationDTO{
		Id:               id,
		Name:             accommodation.Name,
		Location:         accommodation.Location,
		Conveniences:     accommodation.Conveniences,
		MinNumOfVisitors: accommodation.MinNumOfVisitors,
		MaxNumOfVisitors: accommodation.MaxNumOfVisitors,
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
		domainAccommodations = append(domainAccommodations, &domain.AccommodationDTO{
			Id:               id,
			Name:             accommodation.Name,
			Location:         accommodation.Location,
			Conveniences:     accommodation.Conveniences,
			MinNumOfVisitors: accommodation.MinNumOfVisitors,
			MaxNumOfVisitors: accommodation.MaxNumOfVisitors,
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
		Location:         accomm.Location,
		Conveniences:     accomm.Conveniences,
		MinNumOfVisitors: accomm.MinNumOfVisitors,
		MaxNumOfVisitors: accomm.MaxNumOfVisitors,
	}, nil

}

func (as *AccommodationService) UpdateAccommodation(updatedAccommodation domain.Accommodation) (*domain.Accommodation, *errors.ErrorStruct) {
	//as.validator.ValidateAccommodation(&updatedAccommodation)
	log.Println(updatedAccommodation)

	accommodation, _ := as.accommodationRepository.GetAccommodationById(updatedAccommodation.Id.Hex())
	//Check if Name is updated
	if len(updatedAccommodation.Name) != 0 {
		as.validator.ValidateName(updatedAccommodation.Name)
		validatorErrors := as.validator.GetErrors()
		if len(validatorErrors) > 0 {
			var constructedError string
			for _, message := range validatorErrors {
				constructedError += message + "\n"
			}
			return nil, errors.NewError(constructedError, 400)
		}
		accommodation.Name = updatedAccommodation.Name
	}
	//Check if location is updated
	if len(updatedAccommodation.Location) != 0 {
		as.validator.ValidateLocation(updatedAccommodation.Location)
		validatorErrors := as.validator.GetErrors()
		if len(validatorErrors) > 0 {
			var constructedError string
			for _, message := range validatorErrors {
				constructedError += message + "\n"
			}
			return nil, errors.NewError(constructedError, 400)
		}
		accommodation.Location = updatedAccommodation.Location
	}
	//Check if conviniences are updated

	if len(updatedAccommodation.Conveniences) != 0 {
		as.validator.ValidateName(updatedAccommodation.Conveniences)
		validatorErrors := as.validator.GetErrors()
		if len(validatorErrors) > 0 {
			var constructedError string
			for _, message := range validatorErrors {
				constructedError += message + "\n"
			}
			return nil, errors.NewError(constructedError, 400)
		}
		accommodation.Conveniences = updatedAccommodation.Conveniences
	}

	//Check if min number is updated
	if updatedAccommodation.MinNumOfVisitors != 0 {
		log.Println("Uslo je u min num of visitors")
		as.validator.ValidateMinNum(strconv.Itoa(updatedAccommodation.MinNumOfVisitors))
		validatorErrors := as.validator.GetErrors()
		if len(validatorErrors) > 0 {
			var constructedError string
			for _, message := range validatorErrors {
				constructedError += message + "\n"
			}
			return nil, errors.NewError(constructedError, 400)
		}

		accommodation.MinNumOfVisitors = updatedAccommodation.MinNumOfVisitors

	}

	//Check if max number is updated
	if updatedAccommodation.MaxNumOfVisitors != 0 {
		as.validator.ValidateMinNum(strconv.Itoa(updatedAccommodation.MaxNumOfVisitors))
		validatorErrors := as.validator.GetErrors()
		if len(validatorErrors) > 0 {
			var constructedError string
			for _, message := range validatorErrors {
				constructedError += message + "\n"
			}
			return nil, errors.NewError(constructedError, 400)
		}
		accommodation.MaxNumOfVisitors = updatedAccommodation.MaxNumOfVisitors

	}

	if accommodation.MinNumOfVisitors > accommodation.MaxNumOfVisitors {

		return nil, errors.NewError("Min number is higher than max number!", 400)
	}

	log.Println("Prije update")
	_, updateErr := as.accommodationRepository.UpdateAccommodationById(*accommodation)
	if updateErr != nil {
		return nil, errors.NewError("Unable to update", 500)
	}
	log.Println("Poslije update")

	return &domain.Accommodation{
		Id:               updatedAccommodation.Id,
		Name:             updatedAccommodation.Name,
		Location:         updatedAccommodation.Location,
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
