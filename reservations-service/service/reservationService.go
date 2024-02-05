package service

import (
	"context"
	"fmt"
	"log"
	"reservation-service/client"
	"reservation-service/config"
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
	logger       *config.Logger
	tracer       trace.Tracer
}

func NewReservationService(repo *repository.ReservationRepo, validator *utils.Validator, notification *client.NotificationClient, logger *config.Logger, tracer trace.Tracer) *ReservationService {
	return &ReservationService{repo: repo, validator: validator, notification: notification, logger: logger, tracer: tracer}
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
		r.logger.LogError("reservationsService", err.Message)
		return nil, err
	}
	if !available {
		r.logger.LogError("reservationsService", err.Message)
		return nil, errors.NewReservationError(400, "Accommodation not available for the specified date range")
	}

	reserved, erro := r.IsReserved(ctx, reservation.AccommodationID, reservation.DateRange)
	if erro != nil {
		r.logger.LogError("reservationsService", erro.Message)
		return nil, erro
	}
	if reserved {
		r.logger.LogError("reservationsService", erro.Message)
		return nil, errors.NewReservationError(400, "Accommodation not available for the specified date range1")
	}
	createdReservation, insertErr := r.repo.InsertReservation(ctx, &reservation)
	if insertErr != nil {
		r.logger.LogError("reservationsService", erro.Message)
		return nil, errors.NewReservationError(500, "Unable to create reservation: "+insertErr.Error())
	}
	r.notification.SendReservationCreatedNotification(ctx, reservation.HostID, "Reservation successfully created")
	r.logger.LogInfo("reservationsService", fmt.Sprintf("Reservation created: %v", createdReservation))
	return createdReservation, nil
}

func (r ReservationService) CreateAvailability(ctx context.Context, reservation domain.FreeReservation) (*domain.FreeReservation, *errors.ReservationError) {
	ctx, span := r.tracer.Start(ctx, "ReservationService.CreateAvailability")
	defer span.End()

	createdAvailability, insertErr := r.repo.InsertAvailability(ctx, &reservation)
	if insertErr != nil {
		r.logger.LogError("reservationsService", insertErr.Error())
		return nil, errors.NewReservationError(500, "Unable to create availability: "+insertErr.Error())
	}
	r.logger.LogInfo("reservationsService", fmt.Sprintf("Availability created: %v", createdAvailability))
	return createdAvailability, nil
}

func (s *ReservationService) GetReservationsByUser(ctx context.Context, userID string) ([]domain.Reservation, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.GetReservationsByUser")
	defer span.End()

	reservations, err := s.repo.GetReservationsByUser(ctx, userID)
	if err != nil {
		s.logger.LogError("reservationsService", err.Error())

		return nil, errors.NewReservationError(500, err.Error())
	}
	s.logger.LogInfo("reservationsService", fmt.Sprintf("Found reservations by user: %v", reservations))
	return reservations, nil
}
func (s *ReservationService) GetReservationsByHost(ctx context.Context, hostID string) ([]domain.Reservation, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.GetReservationsByHost")
	defer span.End()

	reservations, err := s.repo.GetReservationsByHost(ctx, hostID)
	if err != nil {
		s.logger.LogError("reservationsService", err.Error())
		return nil, errors.NewReservationError(500, err.Error())
	}
	s.logger.LogInfo("reservationsService", fmt.Sprintf("Found reservations by host: %v", reservations))
	return reservations, nil
}

func (s *ReservationService) ProcessDateRange(ctx context.Context, accommodationIDs []string, dateRange []string) ([]string, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.ProcessDateRange")
	defer span.End()
	uniqueAccommodations := make(map[string]struct{})

	reservedAccommodations, err := s.repo.ReservationsInDateRange(ctx, accommodationIDs, dateRange)
	if err != nil {
		s.logger.LogError("reservationsService", err.Error())
		return nil, err
	}

	for _, accommodation := range reservedAccommodations {
		uniqueAccommodations[accommodation] = struct{}{}
	}

	availableAccommodations, err := s.repo.AvailabilityNotInDateRange(ctx, accommodationIDs, dateRange)
	if err != nil {
		s.logger.LogError("reservationsService", err.Error())
		return nil, err
	}

	for _, accommodation := range availableAccommodations {
		uniqueAccommodations[accommodation] = struct{}{}
	}

	result := make([]string, 0, len(uniqueAccommodations))
	for key := range uniqueAccommodations {
		result = append(result, key)
	}
	s.logger.LogInfo("reservationsService", fmt.Sprintf("Found accommodations that have reservations and don't have availability by date range: %v", result))
	return result, nil
}

func (s *ReservationService) GetAvailableDates(ctx context.Context, accommodationID string, dateRange []string) ([]domain.FreeReservation, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.GetAvailableDates")
	defer span.End()
	reservations, err := s.repo.AvailableDates(ctx, accommodationID, dateRange)
	if err != nil {
		s.logger.LogError("reservationsService", err.Error())
		return nil, errors.NewReservationError(500, err.Error())
	}
	s.logger.LogInfo("reservationsService", fmt.Sprintf("Found accommodations availability by date range and accommodationID: %v", reservations))
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
		s.logger.LogError("reservationsService", err.Error())
		return nil, errors.NewReservationError(500, err.Error())
	}
	s.logger.LogInfo("reservationsService", fmt.Sprintf("Found accommodations availability by accommodationID: %v", avl))

	return avl, nil
}

func (s *ReservationService) DeleteReservationById(ctx context.Context, country string, id, userID, hostID, accommodationID, endDate string) (*domain.Reservation, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.DeleteReservationById")
	defer span.End()
	deletedReservation, err := s.repo.DeleteById(ctx, country, id, userID, hostID, accommodationID, endDate)
	if err != nil {
		s.logger.LogError("reservationsService", err.Error())
		return nil, errors.NewReservationError(500, err.Error())
	}
	s.logger.LogInfo("reservationsService", fmt.Sprintf("Deleted reservations by id: %v", deletedReservation))
	return deletedReservation, nil
}

func (s *ReservationService) IsAvailable(ctx context.Context, accommodationID string, dateRange []string) (bool, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.IsAvailable")
	defer span.End()

	available, err := s.repo.IsAvailable(ctx, accommodationID, dateRange)
	if err != nil {
		s.logger.LogError("reservationsService", err.Error())
		return false, errors.NewReservationError(500, "Accommodation not available")
	}
	s.logger.LogInfo("reservationsService", fmt.Sprintf("Found availability for accommodation by accommodationID and date ranges: %v", available))
	return available, nil
}

func (s *ReservationService) IsReserved(ctx context.Context, accommodationID string, dateRange []string) (bool, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.IsReserved")
	defer span.End()

	available, err := s.repo.IsReserved(ctx, accommodationID, dateRange)
	log.Println("available", available)
	if err != nil {
		s.logger.LogError("reservationsService", err.Error())
		return false, errors.NewReservationError(500, "Accommodation not available")
	}
	s.logger.LogInfo("reservationsService", fmt.Sprintf("Checked is accommodation reserved by accommodationID and date ranges: %v", available))
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
		s.logger.LogError("reservationsService", err.Error())
		return 0, errors.NewReservationError(500, "Cannot retrive the number of canceled reservations")
	}
	totalReservations, erro := s.getTotalReservationsByHost(ctx, hostID)
	if erro != nil {
		s.logger.LogError("reservationsService", err.Error())
		return 0, errors.NewReservationError(500, "Cannot retrive the total number of reservations")
	}
	percentageCanceled := float32(numberOfCanceled)/float32(totalReservations) + float32(numberOfCanceled)*100
	s.logger.LogInfo("reservationsService", fmt.Sprintf("Found cancelation percentege: %v", percentageCanceled))
	return percentageCanceled, nil
}

func (s *ReservationService) GetReservationsByAccommodationWithEndDate(ctx context.Context, accommodationID, userID string) ([]domain.Reservation, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.GetReservationsByAccommodationWithEndDate")
	defer span.End()
	reservations, err := s.repo.GetReservationsByAccommodationWithEndDate(ctx, accommodationID, userID)
	if err != nil {
		s.logger.LogError("reservationsService", err.Error())
		return nil, errors.NewReservationError(500, err.Error())
	}
	s.logger.LogInfo("reservationsService", fmt.Sprintf("Found expired reservations by accommodationID: %v", reservations))
	return reservations, nil
}

func (s *ReservationService) GetReservationsByHostWithEndDate(ctx context.Context, hostID, userID string) ([]domain.Reservation, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.GetReservationsByHostWithEndDate")
	defer span.End()
	reservations, err := s.repo.GetReservationsByHostWithEndDate(ctx, hostID, userID)
	if err != nil {
		s.logger.LogError("reservationsService", err.Error())
		return nil, errors.NewReservationError(500, err.Error())
	}
	s.logger.LogInfo("reservationsService", fmt.Sprintf("Found expired reservations by hostID: %v", reservations))
	return reservations, nil
}

func (s *ReservationService) DeleteAvl(ctx context.Context, accommodationID, id, country string, price int) (*domain.FreeReservation, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.DeleteAvl")
	defer span.End()
	deletedAvl, err := s.repo.DeleteAvl(ctx, accommodationID, id, country, price)
	if err != nil {
		s.logger.LogError("reservationsService", err.Error())
		return nil, errors.NewReservationError(500, err.Error())
	}
	s.logger.LogInfo("reservationsService", fmt.Sprintf("Deleted availability by ID: %v", deletedAvl))
	return deletedAvl, nil
}

func (s *ReservationService) UpdateAvailability(ctx context.Context, accommodationID, id, country string, price int, reservation *domain.FreeReservation) (*domain.FreeReservation, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.UpdateAvailability")
	defer span.End()
	_, err := s.repo.DeleteAvl(ctx, accommodationID, id, country, price)
	if err != nil {
		s.logger.LogError("reservationsService", err.Error())
		return nil, errors.NewReservationError(500, err.Error())
	}

	updatedReservation, err := s.repo.InsertAvailability(ctx, reservation)
	if err != nil {
		s.logger.LogError("reservationsService", err.Error())
		return nil, errors.NewReservationError(500, "Unable to update availability!")
	}
	s.logger.LogInfo("reservationsService", fmt.Sprintf("Updated availability by ID: %v", updatedReservation))
	return updatedReservation, nil
}

func (s *ReservationService) GetAccommodationIDsByMaxPrice(ctx context.Context, maxPrice int) ([]string, *errors.ReservationError) {
	ctx, span := s.tracer.Start(ctx, "ReservationService.GetAccommodationIDsByMaxPrice")
	defer span.End()
	accommodations, err := s.repo.GetAccommodationIDsByMaxPrice(ctx, maxPrice)
	if err != nil {
		s.logger.LogError("reservationsService", err.Error())
		return nil, errors.NewReservationError(500, err.Error())
	}
	s.logger.LogInfo("reservationsService", fmt.Sprintf("Found accommodations by price: %v", accommodations))
	return accommodations, nil
}
