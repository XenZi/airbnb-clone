package user_joined

import (
	"encoding/json"
	metrics_events "example/metrics_events"
)

type Event struct {
	UserID                  string `json:"userID"`
	AccommodationID         string `json:"accommodationID"`
	LeftAt                  string `json:"leftAt"`
	number                  uint64
	expectedLastEventNumber int64
}

func NewEvent(userID, accommodationID, leftAt string, expectedLastEventNumber int64) metrics_events.Event {
	return &Event{
		UserID:                  userID,
		AccommodationID:         accommodationID,
		LeftAt:                  leftAt,
		expectedLastEventNumber: expectedLastEventNumber,
	}
}

func NewEmptyEvent() metrics_events.Event {
	return &Event{}
}

func (e *Event) Type() string {
	return metrics_events.EventTypeUserLeft
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
	return "user_left"
}

func (e *Event) ExpectedLastEventNumber() int64 {
	return e.expectedLastEventNumber
}

func (e *Event) SetExpectedLastEventNumber(number uint64) {
	e.expectedLastEventNumber = int64(number)
}
