package handler

import (
	"errors"
	"log"
	"metrics-command/commands"
	"metrics-command/commands/user_joined"
	"metrics-command/commands/user_left"
	"metrics-command/commands/user_rated"
	"metrics-command/commands/user_reserved"
	"metrics-command/store"

	"example/metrics_events"
	user_joined_event "example/metrics_events/user_joined"
	user_left_event "example/metrics_events/user_left"
	user_rated_event "example/metrics_events/user_rated"
	user_reserved_event "example/metrics_events/user_reserved"
)

type Handler struct {
	store store.EventStore
}

func NewHandler(store store.EventStore) Handler {
	return Handler{
		store: store,
	}
}

func (h Handler) Handle(command commands.Command) error {
	event, err := h.execute(command)
	if err != nil {
		return err
	}

	eventJson, err := event.ToJSON()
	if err != nil {
		return err
	}
	log.Println("TYPE EVENT: ", event.Type())
	if event.ExpectedLastEventNumber() >= 0 {
		// we want to enforce event order check
		return h.store.StoreAndExpectLastEventNumber(
			event.Stream(),
			event.Type(),
			eventJson,
			uint64(event.ExpectedLastEventNumber()))
	} else {
		return h.store.Store(event.Stream(), event.Type(), eventJson)
	}
}

func (h Handler) execute(command commands.Command) (event metrics_events.Event, err error) {
	switch c := command.(type) {
	case *user_joined.UserJoinedCommand:
		event, err = h.createUserJoined(c)
	case *user_left.UserLeftCommand:
		event, err = h.createUserLeft(c)
	case *user_reserved.UserReservedCommand:
		event, err = h.createUserReserved(c)
	case *user_rated.UserRatedCommand:
		event, err = h.createUserRated(c)
	default:
		err = errors.New("unknown command")
	}
	return
}

func (h Handler) createUserJoined(command *user_joined.UserJoinedCommand) (metrics_events.Event, error) {
	return user_joined_event.NewEvent(
			command.UserID,
			command.AccommodationID,
			command.JoinedAt,
			command.CustomUUID,
			-1),
		nil
}

func (h Handler) createUserLeft(command *user_left.UserLeftCommand) (metrics_events.Event, error) {
	return user_left_event.NewEvent(
			command.UserID,
			command.AccommodationID,
			command.LeftAt,
			command.CustomUUID,
			-1),
		nil
}

func (h Handler) createUserReserved(command *user_reserved.UserReservedCommand) (metrics_events.Event, error) {
	return user_reserved_event.NewEvent(
			command.UserID,
			command.AccommodationID,
			command.ReservedAt,
			-1),
		nil
}

func (h Handler) createUserRated(command *user_rated.UserRatedCommand) (metrics_events.Event, error) {
	return user_rated_event.NewEvent(
			command.UserID,
			command.AccommodationID,
			command.RatedAt,
			-1),
		nil
}
