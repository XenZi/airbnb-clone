package service

import (
	"reservation-service/domain"
	"reservation-service/errors"
	"reservation-service/repository"
	"reservation-service/utils"
)

type ReservationService struct {
	repo      *repository.ReservationRepo
	validator *utils.Validator
}

func NewReservationService(repo *repository.ReservationRepo, validator *utils.Validator) *ReservationService {
	return &ReservationService{repo: repo, validator: validator}
}

func (r ReservationService) CreateReservationByUser(reservation domain.Reservation) (*domain.Reservation, *errors.ReservationError) {
	r.validator.ValidateReservation(&reservation)
	validatorErrors := r.validator.GetErrors()
	if len(validatorErrors) > 0 {
		var constructedError string
		for _, message := range validatorErrors {
			constructedError += message + "\n"
		}
		return nil, errors.NewReservationError(400, constructedError)
	}
	res, err := r.repo.InsertReservationByUser(&reservation)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}
	return res, nil
}

func (s *ReservationService) GetReservationsByUser(userID string) ([]domain.Reservation, *errors.ReservationError) {

	reservations, err := s.repo.GetReservationsByUser(userID)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}
	return reservations, nil
}

func (s *ReservationService) DeleteReservationById(userId, id string) (*domain.ReservationById, *errors.ReservationError) {
	deletedReservation, err := s.repo.DeleteById(userId, id)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}
	return deletedReservation, nil
}
