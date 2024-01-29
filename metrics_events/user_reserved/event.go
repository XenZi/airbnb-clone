package user_reserved

import (
	"encoding/json"
	metrics_events "metrics-events"
)

type Event struct {
	UserID                  string
	AccommodationID         string
	ReservedAt              string
	expectedLastEventNumber int64
	number                  uint64
}

func NewEvent(userID, accommodationID, reservedAt string, rate int16, expectedLastEventNumber int64) metrics_events.Event {
	return &Event{
		UserID:                  userID,
		AccommodationID:         accommodationID,
		ReservedAt:              reservedAt,
		expectedLastEventNumber: expectedLastEventNumber,
	}
}

func NewEmptyEvent() metrics_events.Event {
	return &Event{}
}

func (e *Event) Type() string {
	return metrics_events.EventTypeUserJoined
}

func (e *Event) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

func (e *Event) FromJSON(jsonEvent []byte) error {
	return json.Unmarshal(jsonEvent, e)
}

func (e *Event) Number() uint64 {
	return e.number
}

func (e *Event) SetNumber(number uint64) {
	e.number = number
}

func (e *Event) Stream() string {
	return "user_joined"
}

func (e *Event) ExpectedLastEventNumber() int64 {
	return e.expectedLastEventNumber
}

func (e *Event) SetExpectedLastEventNumber(number uint64) {
	e.expectedLastEventNumber = int64(number)
}
