package handlers

import (
	"accommodations-service/config"
	"accommodations-service/services"
	events "example/saga/create_accommodation"
	saga "example/saga/messaging"
	"fmt"

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

func (handler *CreateAccommodationCommandHandler) handle(command *events.CreateAccommodationCommand) {
	// ctx, span := handler.tracer.Start(ctx, "CreateAccommodationCommandHandler.handle")
	// defer span.End()
	handler.logger.LogInfo("saga-handler", fmt.Sprintf("USLO U CREATE KOD ACCOMMODATION %v", command.Type))
	returnedValue := command.Payload
	switch command.Type {
	case events.UpdateAccommodation:
		handler.logger.LogInfo("saga-handler", fmt.Sprintf("USLO U CREATE KOD ACCOMMODATION ZA UPDATE ACCOMMODATION %v", command.Type))
		err := handler.accommodationService.ApproveAccommodation(returnedValue.AccommodationID)
		if err != nil {
			return
		}
		break
	case events.DenyAccommodation:
		handler.logger.LogInfo("saga-handler", fmt.Sprintf("USLO U CREATE KOD ACCOMMODATION ZA DENY ACCOMMODATION %v", command.Type))
		err := handler.accommodationService.DenyAccommodation(returnedValue.AccommodationID)
		if err != nil {
			return
		}
		break
	default:
		break
	}
}
