package handlers

import (
	"accommodations-service/services"
	events "example/saga/create_accommodation"
	saga "example/saga/messaging"
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

func (handler *CreateAccommodationCommandHandler) handle(command *events.SagaCommand) {
	returnedID := command.Payload.(string)
	switch command.Type {
	case events.UpdateAccommodation:
		err := handler.accommodationService.ApproveAccommodation(returnedID)
		if err != nil {
			return
		}
	case events.DenyAccommodation:
		err := handler.accommodationService.DenyAccommodation(returnedID)
		if err != nil {
			return
		}
	default:
		return
	}
}
