package orchestrator

import (
	"accommodations-service/config"
	events "example/saga/create_accommodation"
	saga "example/saga/messaging"
	"fmt"
	"log"
)

type CreateAccommodationOrchestrator struct {
	commandPublisher saga.Publisher
	replySubscriber  saga.Subscriber
	logger           *config.Logger
}

func NewCreateAccommodationOrchestrator(publisher saga.Publisher, replySubscriber saga.Subscriber, logger *config.Logger) (*CreateAccommodationOrchestrator, error) {
	orchestrator := &CreateAccommodationOrchestrator{
		commandPublisher: publisher,
		replySubscriber:  replySubscriber,
		logger:           logger,
	}
	err := orchestrator.replySubscriber.Subscribe(orchestrator.handle)
	if err != nil {
		logger.LogError("accommodations-saga-orchestrator", fmt.Sprintf("Unable to subscribe to reply Subscriber"))
		return nil, err
	}
	return orchestrator, nil
}

func (cao *CreateAccommodationOrchestrator) Start(accommodation *events.SendCreateAccommodationAvailability) error {
	log.Println(accommodation.AccommodationID)
	cao.logger.LogInfo("accommodation-saga-orchestrator", "Entered in start saga with id of accommodation "+accommodation.AccommodationID)
	event := &events.CreateAccommodationCommand{
		Type:    events.CreateAvailability,
		Payload: *accommodation,
	}
	cao.logger.LogInfo("accommodation-saga-orchestrator", "Event created ")
	log.Println("EVENT IZGLEDA OVAKO ", event)
	return cao.commandPublisher.Publish(event)
}

func (cao *CreateAccommodationOrchestrator) handle(reply *events.CreateAccommodationReply) {
	cao.logger.LogInfo("accommodation-saga-orchestrator", "Entered saga handle func")
	command := events.CreateAccommodationCommand{Payload: reply.Payload}
	command.Type = cao.nextCommandType(reply.Type)
	if command.Type != events.UnknownCommand {
		cao.logger.LogInfo("accommodation-saga-service", "Entered command publishing")
		_ = cao.commandPublisher.Publish(command)
	}
}

func (cao *CreateAccommodationOrchestrator) nextCommandType(reply events.CreateAccommodationReplyType) events.CreateAccommodationCommandType {
	cao.logger.LogInfo("accommodation-saga-service", "Picking next command")
	switch reply {
	case events.AvailabilityCreated:
		return events.UpdateAccommodation
	case events.AvailabilityNotCreated:
		return events.RollbackAccommodation
	default:
		return events.UnknownCommand
	}

}
