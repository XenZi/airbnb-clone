package orchestrator

import (
	"accommodations-service/domain"
	events "example/saga/create_accommodation"
	saga "example/saga/messaging"
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

func (cao *CreateAccommodationOrchestrator) Start(accommodation *domain.SendCreateAccommodationAvailability) error {
	event := &events.SagaCommand{
		Type:    events.CreateAvailability,
		Payload: accommodation,
	}
	return cao.commandPublisher.Publish(event)
}

func (cao *CreateAccommodationOrchestrator) handle(reply *events.SagaReply) {
	command := events.SagaCommand{Payload: reply.Payload}
	command.Type = cao.nextCommandType(reply.Type)
	if command.Type != events.UnknownCommand {
		_ = cao.commandPublisher.Publish(command)
	}
}

func (cao *CreateAccommodationOrchestrator) nextCommandType(reply events.SagaReplyType) events.SagaCommandType {
	switch reply {
	case events.AvailabilityCreated:
		return events.UpdateAccommodation
	case events.AvailabilityNotCreated:
		return events.RollbackAccommodation
	default:
		return events.UnknownCommand
	}

}
