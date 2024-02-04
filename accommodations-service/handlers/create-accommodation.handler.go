package handlers

import (
	"accommodations-service/services"
	"context"
	events "example/saga/create_accommodation"
	saga "example/saga/messaging"
	"log"

	"go.opentelemetry.io/otel/trace"
)

type CreateAccommodationCommandHandler struct {
	accommodationService *services.AccommodationService
	replyPublisher       saga.Publisher
	commandSubscriber    saga.Subscriber
	tracer               trace.Tracer
}

func NewCreateAccommodationCommandHandler(accommodationService *services.AccommodationService, publisher saga.Publisher, subscriber saga.Subscriber, tracer trace.Tracer) (*CreateAccommodationCommandHandler, error) {
	o := &CreateAccommodationCommandHandler{
		accommodationService: accommodationService,
		replyPublisher:       publisher,
		commandSubscriber:    subscriber,
		tracer:               tracer,
	}
	err := o.commandSubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (handler *CreateAccommodationCommandHandler) handle(ctx context.Context, command *events.CreateAccommodationCommand) {
	ctx, span := handler.tracer.Start(ctx, "CreateAccommodationCommandHandler.handle")
	defer span.End()
	log.Println("KOMANDA USLA U CREATE AVAILABILITY KOD ACCOMMODATIONS SERVICE", command.Type)
	returnedValue := command.Payload
	switch command.Type {
	case events.UpdateAccommodation:
		err := handler.accommodationService.ApproveAccommodation(ctx, returnedValue.AccommodationID)
		if err != nil {
			return
		}
		break
	case events.DenyAccommodation:
		err := handler.accommodationService.DenyAccommodation(ctx, returnedValue.AccommodationID)
		if err != nil {
			return
		}
		break
	default:
		break
	}
}
