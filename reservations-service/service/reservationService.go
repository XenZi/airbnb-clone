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

// service/reservationService.go

func (r ReservationService) CreateReservation(reservation domain.Reservation) (*domain.Reservation, *errors.ReservationError) {
	r.validator.ValidateReservation(&reservation)
	validationErrors := r.validator.GetErrors()

	if len(validationErrors) > 0 {
		return nil, errors.NewReservationError(400, "Validation failed")
	}

	/*	available, err := r.AvailableDates(reservation.AccommodationID, reservation.StartDate, reservation.EndDate)
		if err != nil {
			return nil, err
		}

			if !available {
			return nil, errors.NewReservationError(400, "Accommodation not available for the specified date range")
		} */
	createdReservation, insertErr := r.repo.InsertReservation(&reservation)
	if insertErr != nil {
		return nil, errors.NewReservationError(500, "Unable to create reservation: "+insertErr.Error())
	}

	return createdReservation, nil
}

func (r ReservationService) CreateAvailability(reservation domain.FreeReservation) (*domain.FreeReservation, *errors.ReservationError) {
	r.validator.ValidateAvailability(&reservation)
	validationErrors := r.validator.GetErrors()

	if len(validationErrors) > 0 {
		return nil, errors.NewReservationError(400, "Validation failed")
	}
	createdAvailability, insertErr := r.repo.InsertAvailability(&reservation)
	if insertErr != nil {
		return nil, errors.NewReservationError(500, "Unable to create availability: "+insertErr.Error())
	}
	return createdAvailability, nil
}

func (s *ReservationService) GetReservationsByUser(userID string) ([]domain.Reservation, *errors.ReservationError) {

	reservations, err := s.repo.GetReservationsByUser(userID)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}
	return reservations, nil
}

func (s *ReservationService) GetReservationsByAccommodation(accommodationID string) ([]domain.Reservation, *errors.ReservationError) {
	reservations, err := s.repo.GetReservationsByAccommodation(accommodationID)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}
	return reservations, nil
}

func (s *ReservationService) DeleteReservationById(id string) (*domain.ReservationById, *errors.ReservationError) {
	deletedReservation, err := s.repo.DeleteById(id)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}
	return deletedReservation, nil
}
func (s *ReservationService) AvailableDates(accommodationID string, startDate, endDate string) (bool, *errors.ReservationError) {

	available, err := s.repo.AvailableDates(accommodationID, startDate, endDate)
	if err != nil {
		return false, errors.NewReservationError(400, "Accommodation not available")
	}

	return available, nil
}
