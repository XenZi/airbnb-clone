package services

import (
	"github.com/XenZi/airbnb-clone/accommodations-service/domain"
	"github.com/XenZi/airbnb-clone/accommodations-service/errors"
	events "github.com/XenZi/airbnb-clone/saga/create_accommodation"
	saga "github.com/XenZi/airbnb-clone/saga/messaging"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (o *CreateAccommodationOrchestrator) handle(reply *events.CreateAccommodationReply) {
	command := events.CreateAccommodationCommand{Accommodation: reply.Accommodation}
	command.Type = o.nextCommandType(reply.Type)
	if command.Type != events.UnknownCommand {
		_ = o.commandPublisher.Publish(command)
	}
}

func (o *CreateAccommodationOrchestrator) Start(accommodation domain.CreateAccommodation, id string) error {
	idPrimitive, _ := primitive.ObjectIDFromHex(id)

	event := &events.CreateAccommodationCommand{
		Type: events.UpdateAccommodation,
		Accommodation: events.CreateAccommodation{
			Id:               idPrimitive,
			Name:             accommodation.Name,
			Address:          accommodation.Address,
			City:             accommodation.City,
			Country:          accommodation.Country,
			UserName:         accommodation.UserName,
			UserId:           accommodation.UserId,
			Email:            accommodation.Email,
			Conveniences:     accommodation.Conveniences,
			MinNumOfVisitors: accommodation.MinNumOfVisitors,
			MaxNumOfVisitors: accommodation.MaxNumOfVisitors,
			Status:           events.AccommodationStatus(domain.Pending),
		},
	}

	return o.commandPublisher.Publish(event)
}

func (o *CreateAccommodationOrchestrator) nextCommandType(reply events.CreateAccommodationReplyType) events.CreateAccommodationCommandType {
	switch reply {
	case events.AccommodationCreated:
		return events.UpdateAccommodation
	case events.AccommodationNotCreated:
		return events.DeleteAccommodation
	case events.AvailabilitiesCreated:
		return events.ApproveAccommodation
	case events.AvailabilitiesNotCreated:
		return events.DeleteAccommodation
	default:
		return events.UnknownCommand
	}
}
