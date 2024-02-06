package handlers

import (
	"accommodations-service/config"
	"accommodations-service/services"
	"context"
	events "example/saga/create_accommodation"
	saga "example/saga/messaging"
	"fmt"
	"log"

	"go.opentelemetry.io/otel/trace"
)

type CreateAccommodationCommandHandler struct {
	accommodationService *services.AccommodationService
	replyPublisher       saga.Publisher
	commandSubscriber    saga.Subscriber
	tracer               trace.Tracer
	logger               *config.Logger
}

func NewCreateAccommodationCommandHandler(accommodationService *services.AccommodationService, publisher saga.Publisher, subscriber saga.Subscriber, tracer trace.Tracer, logger *config.Logger) (*CreateAccommodationCommandHandler, error) {
	o := &CreateAccommodationCommandHandler{
		accommodationService: accommodationService,
		replyPublisher:       publisher,
		commandSubscriber:    subscriber,
		tracer:               tracer,
		logger:               logger,
	}
	err := o.commandSubscriber.Subscribe(o.handle)
	if err != nil {
		logger.LogError("accommodations-saga-handler", fmt.Sprintf("Unable to subscribe to commandSubscriber"))
		return nil, err
	}
	return o, nil
}

func (handler *CreateAccommodationCommandHandler) handle(ctx context.Context, command *events.CreateAccommodationCommand) {
	ctx, span := handler.tracer.Start(ctx, "CreateAccommodationCommandHandler.handle")
	defer span.End()
	handler.logger.LogInfo("accommodation-saga-handler", "Entered in create availability in handler")
	log.Println("KOMANDA USLA U CREATE AVAILABILITY KOD ACCOMMODATIONS SERVICE", command.Type)
	returnedValue := command.Payload
	switch command.Type {
	case events.UpdateAccommodation:
		handler.logger.LogInfo("accommodation-saga-handler", "Entered updating accommodation")
		err := handler.accommodationService.ApproveAccommodation(ctx, returnedValue.AccommodationID)
		if err != nil {
			return
		}
		break
	case events.DenyAccommodation:
		handler.logger.LogInfo("accommodation-saga-handler", "Entered Denying accommodation accommodation")
		err := handler.accommodationService.DenyAccommodation(ctx, returnedValue.AccommodationID)
		if err != nil {
			return
		}
		break
	default:
		break
	}
}
