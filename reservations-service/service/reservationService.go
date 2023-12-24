package service

import (
	client "command-line-arguments/home/janko33/Documents/airbnb-clone/accommodations-service/client/reservations-client.client.go"
	"context"
	"log"
	"reservation-service/domain"
	"reservation-service/errors"
	"reservation-service/repository"
	"reservation-service/utils"
)

type ReservationService struct {
	repo         *repository.ReservationRepo
	validator    *utils.Validator
	notification *client.NotificationClient
}

func NewReservationService(repo *repository.ReservationRepo, validator *utils.Validator, notification *client.NotificationClient) *ReservationService {
	return &ReservationService{repo: repo, validator: validator, notification: notification}
}

// service/reservationService.go

func (r ReservationService) CreateReservation(reservation domain.Reservation, ctx context.Context) (*domain.Reservation, *errors.ReservationError) {
	r.validator.ValidateReservation(&reservation)
	validationErrors := r.validator.GetErrors()

	if len(validationErrors) > 0 {
		return nil, errors.NewReservationError(400, "Validation failed")
	}
	log.Println(reservation.StartDate, reservation.EndDate)
	available, err := r.IsAvailable(reservation.AccommodationID, reservation.StartDate, reservation.EndDate)
	if err != nil {
		return nil, err
	}

	if !available {
		return nil, errors.NewReservationError(400, "Accommodation not available for the specified date range")
	}
	reserved, erro := r.IsReserved(reservation.AccommodationID, reservation.StartDate, reservation.EndDate)
	if erro != nil {
		return nil, erro
	}
	if reserved {
		return nil, errors.NewReservationError(400, "Accommodation not available for the specified date range1")
	}
	createdReservation, insertErr := r.repo.InsertReservation(&reservation)
	if insertErr != nil {
		return nil, errors.NewReservationError(500, "Unable to create reservation: "+insertErr.Error())
	}
	r.notification.SendReservationCreatedNotification(ctx, reservation.UserID, "Reservation successfully created")
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
func (s *ReservationService) GetReservationsByHost(hostID string) ([]domain.Reservation, *errors.ReservationError) {

	reservations, err := s.repo.GetReservationsByUser(hostID)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}
	return reservations, nil
}
func (s *ReservationService) ReservationsInDateRange(accommodationIDs []string, dateRange []string) ([]string, *errors.ReservationError) {
	reservations, err := s.repo.ReservationsInDateRange(accommodationIDs, dateRange)
	if err != nil {
		return nil, err
	}
	return reservations, nil
}
func (s *ReservationService) GetAvailableDates(accommodationID, startDate, endDate string) ([]domain.FreeReservation, *errors.ReservationError) {
	reservations, err := s.repo.AvailableDates(accommodationID, startDate, endDate)
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
func (s *ReservationService) GetAvailabilityForAccommodation(accommodationID string) ([]domain.GetAvailabilityForAccommodation, *errors.ReservationError) {
	avl, err := s.repo.CheckAvailabilityForAccommodation(accommodationID)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}
	return avl, nil
}

func (s *ReservationService) DeleteReservationById(country string, id string) (*domain.Reservation, *errors.ReservationError) {
	deletedReservation, err := s.repo.DeleteById(country, id)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}
	return deletedReservation, nil
}

func (s *ReservationService) IsAvailable(accommodationID string, startDate, endDate string) (bool, *errors.ReservationError) {

	available, err := s.repo.IsAvailable(accommodationID, startDate, endDate)
	if err != nil {
		return false, errors.NewReservationError(400, "Accommodation not available")
	}

	return available, nil
}
func (s *ReservationService) IsReserved(accommodationID string, startDate, endDate string) (bool, *errors.ReservationError) {

	available, err := s.repo.IsReserved(accommodationID, startDate, endDate)
	if err != nil {
		return false, errors.NewReservationError(400, "Accommodation not available")
	}

	return available, nil
}
