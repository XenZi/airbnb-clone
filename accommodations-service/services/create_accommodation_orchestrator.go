package services

import (
	"github.com/XenZi/airbnb-clone/accommodations-service/errors"
	saga "github.com/XenZi/airbnb-clone/saga/messaging"
)

type CreateAccommodationOrchestrator struct {
	commandPublisher saga.Publisher
	replySubscriber  saga.Subscriber
}

func NewCreateAccommodationOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) (*CreateAccommodationOrchestrator, *errors.ErrorStruct) {
	o := &CreateAccommodationOrchestrator{
		commandPublisher: publisher,
		replySubscriber:  subscriber,
	}
	err := o.replySubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, errors.NewError("Not subscribed correctly", 500)
	}
	return o, nil
}

func (o *CreateOrderOrchestrator) handle(reply *events.CreateOrderReply) {
	command := events.CreateOrderCommand{Order: reply.Order}
	command.Type = o.nextCommandType(reply.Type)
	if command.Type != events.UnknownCommand {
		_ = o.commandPublisher.Publish(command)
	}
}
