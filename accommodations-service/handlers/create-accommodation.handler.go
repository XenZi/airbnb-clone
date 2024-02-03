package handlers

import (
	"accommodations-service/services"
	events "example/saga/create_accommodation"
	saga "example/saga/messaging"
	"log"
)

type CreateAccommodationCommandHandler struct {
	accommodationService *services.AccommodationService
	replyPublisher       saga.Publisher
	commandSubscriber    saga.Subscriber
}

func NewCreateAccommodationCommandHandler(accommodationService *services.AccommodationService, publisher saga.Publisher, subscriber saga.Subscriber) (*CreateAccommodationCommandHandler, error) {
	o := &CreateAccommodationCommandHandler{
		accommodationService: accommodationService,
		replyPublisher:       publisher,
		commandSubscriber:    subscriber,
	}
	err := o.commandSubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (handler *CreateAccommodationCommandHandler) handle(command *events.CreateAccommodationCommand) {
	log.Println("KOMANDA USLA U CREATE AVAILABILITY KOD ACCOMMODATIONS SERVICE", command.Type)
	returnedValue := command.Payload
	switch command.Type {
	case events.UpdateAccommodation:
		err := handler.accommodationService.ApproveAccommodation(returnedValue.AccommodationID)
		if err != nil {
			return
		}
		break
	case events.DenyAccommodation:
		err := handler.accommodationService.DenyAccommodation(returnedValue.AccommodationID)
		if err != nil {
			return
		}
		break
	default:
		break
	}
}
