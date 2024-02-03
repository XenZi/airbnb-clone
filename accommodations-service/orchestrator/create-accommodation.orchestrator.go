package orchestrator

import (
	events "example/saga/create_accommodation"
	saga "example/saga/messaging"
	"log"
)

type CreateAccommodationOrchestrator struct {
	commandPublisher saga.Publisher
	replySubscriber  saga.Subscriber
}

func NewCreateAccommodationOrchestrator(publisher saga.Publisher, replySubscriber saga.Subscriber) (*CreateAccommodationOrchestrator, error) {
	orchestrator := &CreateAccommodationOrchestrator{
		commandPublisher: publisher,
		replySubscriber:  replySubscriber,
	}
	err := orchestrator.replySubscriber.Subscribe(orchestrator.handle)
	if err != nil {
		return nil, err
	}
	return orchestrator, nil
}

func (cao *CreateAccommodationOrchestrator) Start(accommodation *events.SendCreateAccommodationAvailability) error {
	log.Println(accommodation.AccommodationID)
	event := &events.CreateAccommodationCommand{
		Type:    events.CreateAvailability,
		Payload: *accommodation,
	}
	log.Println("EVENT IZGLEDA OVAKO ", event)
	return cao.commandPublisher.Publish(event)
}

func (cao *CreateAccommodationOrchestrator) handle(reply *events.CreateAccommodationReply) {
	command := events.CreateAccommodationCommand{Payload: reply.Payload}
	command.Type = cao.nextCommandType(reply.Type)
	if command.Type != events.UnknownCommand {
		_ = cao.commandPublisher.Publish(command)
	}
}

func (cao *CreateAccommodationOrchestrator) nextCommandType(reply events.CreateAccommodationReplyType) events.CreateAccommodationCommandType {
	switch reply {
	case events.AvailabilityCreated:
		return events.UpdateAccommodation
	case events.AvailabilityNotCreated:
		return events.RollbackAccommodation
	default:
		return events.UnknownCommand
	}

}
