package service

import (
	"context"
	"log"
	"reservation-service/client"
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

	available, err := r.IsAvailable(reservation.AccommodationID, reservation.DateRange)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, errors.NewReservationError(400, "Accommodation not available for the specified date range")
	}

	reserved, erro := r.IsReserved(reservation.AccommodationID, reservation.DateRange)
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
	r.notification.SendReservationCreatedNotification(ctx, reservation.HostID, "Reservation successfully created")
	return createdReservation, nil
}

func (r ReservationService) CreateAvailability(reservation domain.FreeReservation) (*domain.FreeReservation, *errors.ReservationError) {

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
func (s *ReservationService) GetReservationsByAccommodationWithEndDate(accommodationID, userID string) ([]domain.Reservation, *errors.ReservationError) {
	reservations, err := s.repo.GetReservationsByAccommodationWithEndDate(accommodationID, userID)
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
func (s *ReservationService) GetAvailableDates(accommodationID string, dateRange []string) ([]domain.FreeReservation, *errors.ReservationError) {
	reservations, err := s.repo.AvailableDates(accommodationID, dateRange)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}
	return reservations, nil
}

/*
	func (s *ReservationService) GetReservationsByAccommodation(accommodationID string) ([]domain.Reservation, *errors.ReservationError) {
		reservations, err := s.repo.GetReservationsByAccommodation(accommodationID)
		if err != nil {
			return nil, errors.NewReservationError(500, err.Error())
		}
		return reservations, nil
	}
*/
func (s *ReservationService) GetAvailabilityForAccommodation(accommodationID string) ([]domain.GetAvailabilityForAccommodation, *errors.ReservationError) {
	avl, err := s.repo.CheckAvailabilityForAccommodation(accommodationID)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}
	return avl, nil
}

func (s *ReservationService) DeleteReservationById(country string, id, userID, hostID, accommodationID, endDate string) (*domain.Reservation, *errors.ReservationError) {
	deletedReservation, err := s.repo.DeleteById(country, id, userID, hostID, accommodationID, endDate)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}
	return deletedReservation, nil
}

func (s *ReservationService) IsAvailable(accommodationID string, dateRange []string) (bool, *errors.ReservationError) {

	available, err := s.repo.IsAvailable(accommodationID, dateRange)
	if err != nil {
		return false, errors.NewReservationError(500, "Accommodation not available")
	}

	return available, nil
}
func (s *ReservationService) IsReserved(accommodationID string, dateRange []string) (bool, *errors.ReservationError) {

	available, err := s.repo.IsReserved(accommodationID, dateRange)
	log.Println("available", available)
	if err != nil {
		return false, errors.NewReservationError(500, "Accommodation not available")
	}

	return available, nil
}
func (s *ReservationService) getNumberOfCanceledReservations(hostID string) (int, *errors.ReservationError) {
	numberOfCanceledReservations, err := s.repo.GetNumberOfCanceledReservations(hostID)
	if err != nil {
		return 0, errors.NewReservationError(500, "Cannot retrive the number of canceled reservations")
	}

	return numberOfCanceledReservations, nil
}
func (s *ReservationService) getTotalReservationsByHost(hostID string) (int, *errors.ReservationError) {
	totalReservations, err := s.repo.GetTotalReservationsByHost(hostID)
	if err != nil {
		return 0, errors.NewReservationError(500, "Cannot retrive the total number of reservations")
	}
	return totalReservations, nil
}

func (s *ReservationService) CalculatePercentageCanceled(hostID string) (float32, *errors.ReservationError) {
	numberOfCanceled, err := s.getNumberOfCanceledReservations(hostID)
	if err != nil {
		return 0, errors.NewReservationError(500, "Cannot retrive the number of canceled reservations")
	}
	totalReservations, erro := s.getTotalReservationsByHost(hostID)
	if erro != nil {
		return 0, errors.NewReservationError(500, "Cannot retrive the total number of reservations")
	}
	percentageCanceled := float32(numberOfCanceled)/float32(totalReservations) + float32(numberOfCanceled)*100
	return percentageCanceled, nil
}
