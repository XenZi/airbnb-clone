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

func (r ReservationService) CreateReservationByUser(reservation domain.Reservation) (*domain.Reservation, *errors.ReservationError) {
	r.validator.ValidateReservation(&reservation)
	validationErrors := r.validator.GetErrors()

	if len(validationErrors) > 0 {
		return nil, errors.NewReservationError(400, "Validation failed")
	}

	available, err := r.CheckAvailability(reservation.AccommodationID, reservation.StartDate, reservation.EndDate)
	if err != nil {
		return nil, err
	}

	if !available {
		return nil, errors.NewReservationError(400, "Accommodation not available for the specified date range")
	}
	createdReservation, insertErr := r.repo.InsertReservationByUser(&reservation)
	if insertErr != nil {
		return nil, errors.NewReservationError(500, "Unable to create reservation: "+insertErr.Error())
	}

	return createdReservation, nil
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
func (s *ReservationService) CheckAvailability(accommodationID string, startDate, endDate string) (bool, *errors.ReservationError) {

	available, err := s.repo.CheckAvailability(accommodationID, startDate, endDate)
	if err != nil {
		return false, errors.NewReservationError(400, "Accommodation not available")
	}

	return available, nil
}
