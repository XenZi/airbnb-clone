package services

import (
	"accommodations-service/errors"
	saga "saga/messaging"
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
