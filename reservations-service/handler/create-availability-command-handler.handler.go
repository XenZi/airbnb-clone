package handler

import (
	"context"
	events "example/saga/create_accommodation"
	saga "example/saga/messaging"
	"log"
	"reservation-service/domain"
	"reservation-service/service"
)

type CreateAvailabilityCommandHandler struct {
	reservationService *service.ReservationService
	replyPublisher     saga.Publisher
	commandSubscriber  saga.Subscriber
}

type SendCreateAccommodationAvailability struct {
	AccommodationID string                        `json:"accommodationId"`
	Location        string                        `json:"location"`
	DateRange       []AvailableAccommodationDates `json:"dateRange"`
}

type AvailableAccommodationDates struct {
	AccommodationId string   `json:"accommodationId"`
	DateRange       []string `json:"dateRange"`
	Location        string   `json:"location"`
	Price           int      `json:"price"`
}

func NewCreateAvailabilityCommandHandler(reservationService *service.ReservationService, replyPublisher saga.Publisher, commandSubscriber saga.Subscriber) (*CreateAvailabilityCommandHandler, error) {
	o := &CreateAvailabilityCommandHandler{
		reservationService: reservationService,
		replyPublisher:     replyPublisher,
		commandSubscriber:  commandSubscriber,
	}
	err := o.commandSubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (handler CreateAvailabilityCommandHandler) handle(command *events.CreateAccommodationCommand) {
	log.Println("KOMANDA USLA U CREATE AVAILABILITY", command.Type)
	valueFromCommand := command.Payload
	reply := events.CreateAccommodationReply{Payload: valueFromCommand}
	switch command.Type {
	case events.CreateAvailability:
		var dateRangeCasted []domain.DateRangeWithPrice
		for _, value := range valueFromCommand.DateRange {
			val := domain.DateRangeWithPrice{
				DateRange: value.DateRange,
				Price:     value.Price,
			}
			dateRangeCasted = append(dateRangeCasted, val)
		}
		freeAccommodation := domain.FreeReservation{
			AccommodationID: valueFromCommand.AccommodationID,
			Location:        valueFromCommand.Location,
			DateRange:       dateRangeCasted,
		}

		_, err := handler.reservationService.CreateAvailability(context.Background(), freeAccommodation)
		if err != nil {
			reply.Type = events.AvailabilityNotCreated
			break
		}
		reply.Type = events.AvailabilityCreated
		break
	default:
		reply.Type = events.UnknownReply
		break
	}
	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
