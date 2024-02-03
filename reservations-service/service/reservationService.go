package service

import (
	"context"
	"log"
	"reservation-service/client"
	"reservation-service/domain"
	"reservation-service/errors"
	"reservation-service/repository"
	"reservation-service/utils"

	"go.opentelemetry.io/otel/trace"
)

type ReservationService struct {
	repo         *repository.ReservationRepo
	validator    *utils.Validator
	notification *client.NotificationClient
	tracer       trace.Tracer
}

func NewReservationService(repo *repository.ReservationRepo, validator *utils.Validator, notification *client.NotificationClient, tracer trace.Tracer) *ReservationService {
	return &ReservationService{repo: repo, validator: validator, notification: notification, tracer: tracer}
}

// service/reservationService.go

func (r ReservationService) CreateReservation(ctx context.Context, reservation domain.Reservation) (*domain.Reservation, *errors.ReservationError) {
	ctx, span := r.tracer.Start(ctx, "ReservationService.CreateReservation")
	defer span.End()
	r.validator.ValidateReservation(&reservation)
	validationErrors := r.validator.GetErrors()

	if len(validationErrors) > 0 {
		return nil, errors.NewReservationError(400, "Validation failed")
	}

	available, err := r.IsAvailable(ctx, reservation.AccommodationID, reservation.DateRange)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, errors.NewReservationError(400, "Accommodation not available for the specified date range")
	}

	reserved, erro := r.IsReserved(ctx, reservation.AccommodationID, reservation.DateRange)
	if erro != nil {
		return nil, erro
	}
	if reserved {
		return nil, errors.NewReservationError(400, "Accommodation not available for the specified date range1")
	}
	createdReservation, insertErr := r.repo.InsertReservation(ctx, &reservation)
	if insertErr != nil {
		return nil, errors.NewReservationError(500, "Unable to create reservation: "+insertErr.Error())
	}
	r.notification.SendReservationCreatedNotification(ctx, reservation.HostID, "Reservation successfully created")
	return createdReservation, nil
}

func (r ReservationService) CreateAvailability(ctx context.Context, reservation domain.FreeReservation) (*domain.FreeReservation, *errors.ReservationError) {
	ctx, span := r.tracer.Start(ctx, "ReservationService.CreateAvailability")
	defer span.End()

	createdAvailability, insertErr := r.repo.InsertAvailability(ctx, &reservation)
	if insertErr != nil {
		return nil, errors.NewReservationError(500, "Unable to create availability: "+insertErr.Error())
	}
	return createdAvailability, nil
}

func (s *ReservationService) GetReservationsByUser(ctx context.Context, userID string) ([]domain.Reservation, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.GetReservationsByUser")
	defer span.End()

	reservations, err := s.repo.GetReservationsByUser(ctx, userID)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}
	return reservations, nil
}
func (s *ReservationService) GetReservationsByHost(ctx context.Context, hostID string) ([]domain.Reservation, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.GetReservationsByHost")
	defer span.End()

	reservations, err := s.repo.GetReservationsByUser(ctx, hostID)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}
	return reservations, nil
}

func (s *ReservationService) ProcessDateRange(ctx context.Context, accommodationIDs []string, dateRange []string) ([]string, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.ProcessDateRange")
	defer span.End()
	uniqueAccommodations := make(map[string]struct{})

	reservedAccommodations, err := s.repo.ReservationsInDateRange(ctx, accommodationIDs, dateRange)
	if err != nil {
		return nil, err
	}

	for _, accommodation := range reservedAccommodations {
		uniqueAccommodations[accommodation] = struct{}{}
	}

	availableAccommodations, err := s.repo.AvailabilityNotInDateRange(ctx, accommodationIDs, dateRange)
	if err != nil {
		return nil, err
	}

	for _, accommodation := range availableAccommodations {
		uniqueAccommodations[accommodation] = struct{}{}
	}

	result := make([]string, 0, len(uniqueAccommodations))
	for key := range uniqueAccommodations {
		result = append(result, key)
	}

	return result, nil
}

func (s *ReservationService) GetAvailableDates(ctx context.Context, accommodationID string, dateRange []string) ([]domain.FreeReservation, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.GetAvailableDates")
	defer span.End()
	reservations, err := s.repo.AvailableDates(ctx, accommodationID, dateRange)
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
func (s *ReservationService) GetAvailabilityForAccommodation(ctx context.Context, accommodationID string) ([]domain.GetAvailabilityForAccommodation, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.GetAvailabilityForAccommodation")
	defer span.End()
	avl, err := s.repo.CheckAvailabilityForAccommodation(ctx, accommodationID)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}
	return avl, nil
}

func (s *ReservationService) DeleteReservationById(ctx context.Context, country string, id, userID, hostID, accommodationID, endDate string) (*domain.Reservation, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.DeleteReservationById")
	defer span.End()
	deletedReservation, err := s.repo.DeleteById(ctx, country, id, userID, hostID, accommodationID, endDate)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}
	return deletedReservation, nil
}

func (s *ReservationService) IsAvailable(ctx context.Context, accommodationID string, dateRange []string) (bool, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.IsAvailable")
	defer span.End()

	available, err := s.repo.IsAvailable(ctx, accommodationID, dateRange)
	if err != nil {
		return false, errors.NewReservationError(500, "Accommodation not available")
	}

	return available, nil
}

func (s *ReservationService) IsReserved(ctx context.Context, accommodationID string, dateRange []string) (bool, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.IsReserved")
	defer span.End()

	available, err := s.repo.IsReserved(ctx, accommodationID, dateRange)
	log.Println("available", available)
	if err != nil {
		return false, errors.NewReservationError(500, "Accommodation not available")
	}

	return available, nil
}

func (s *ReservationService) getNumberOfCanceledReservations(ctx context.Context, hostID string) (int, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.getNumberOfCanceledReservations")
	defer span.End()
	numberOfCanceledReservations, err := s.repo.GetNumberOfCanceledReservations(ctx, hostID)
	if err != nil {
		return 0, errors.NewReservationError(500, "Cannot retrive the number of canceled reservations")
	}

	return numberOfCanceledReservations, nil
}

func (s *ReservationService) getTotalReservationsByHost(ctx context.Context, hostID string) (int, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.getTotalReservationsByHost")
	defer span.End()
	totalReservations, err := s.repo.GetTotalReservationsByHost(ctx, hostID)
	if err != nil {
		return 0, errors.NewReservationError(500, "Cannot retrive the total number of reservations")
	}
	return totalReservations, nil
}

func (s *ReservationService) CalculatePercentageCanceled(ctx context.Context, hostID string) (float32, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.CalculatePercentageCanceled")
	defer span.End()
	numberOfCanceled, err := s.getNumberOfCanceledReservations(ctx, hostID)
	if err != nil {
		return 0, errors.NewReservationError(500, "Cannot retrive the number of canceled reservations")
	}
	totalReservations, erro := s.getTotalReservationsByHost(ctx, hostID)
	if erro != nil {
		return 0, errors.NewReservationError(500, "Cannot retrive the total number of reservations")
	}
	percentageCanceled := float32(numberOfCanceled)/float32(totalReservations) + float32(numberOfCanceled)*100
	return percentageCanceled, nil
}

func (s *ReservationService) GetReservationsByAccommodationWithEndDate(ctx context.Context, accommodationID, userID string) ([]domain.Reservation, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.GetReservationsByAccommodationWithEndDate")
	defer span.End()
	reservations, err := s.repo.GetReservationsByAccommodationWithEndDate(ctx, accommodationID, userID)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}
	return reservations, nil
}

func (s *ReservationService) GetReservationsByHostWithEndDate(ctx context.Context, hostID, userID string) ([]domain.Reservation, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.GetReservationsByHostWithEndDate")
	defer span.End()
	reservations, err := s.repo.GetReservationsByHostWithEndDate(ctx, hostID, userID)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}
	return reservations, nil
}

func (s *ReservationService) DeleteAvl(ctx context.Context, accommodationID, id, country string, price int) (*domain.FreeReservation, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.DeleteAvl")
	defer span.End()
	deletedAvl, err := s.repo.DeleteAvl(ctx, accommodationID, id, country, price)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}
	return deletedAvl, nil
}

func (s *ReservationService) UpdateAvailability(ctx context.Context, accommodationID, id, country string, price int, reservation *domain.FreeReservation) (*domain.FreeReservation, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.UpdateAvailability")
	defer span.End()
	_, err := s.repo.DeleteAvl(ctx, accommodationID, id, country, price)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}

	updatedReservation, err := s.repo.InsertAvailability(ctx, reservation)
	if err != nil {
		return nil, errors.NewReservationError(500, "Unable to update availability!")
	}

	return updatedReservation, nil
}

func (s *ReservationService) GetAccommodationIDsByMaxPrice(ctx context.Context, maxPrice int) ([]string, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.GetAccommodationIDsByMaxPrice")
	defer span.End()
	accommodations, err := s.repo.GetAccommodationIDsByMaxPrice(ctx, maxPrice)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}

	return accommodations, nil
}
