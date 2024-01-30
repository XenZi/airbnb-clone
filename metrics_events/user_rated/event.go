package user_rated

import (
	"encoding/json"
	metrics_events "example/metrics_events"
)

type Event struct {
	UserID                  string
	AccommodationID         string
	RatedAt                 string
	expectedLastEventNumber int64
	number                  uint64
}

func NewEvent(userID, accommodationID, ratedAt string, expectedLastEventNumber int64) metrics_events.Event {
	return &Event{
		UserID:                  userID,
		AccommodationID:         accommodationID,
		RatedAt:                 ratedAt,
		expectedLastEventNumber: expectedLastEventNumber,
	}
}

func NewEmptyEvent() metrics_events.Event {
	return &Event{}
}

func (e *Event) Type() string {
	return metrics_events.EventTypeUserRated
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
	return "user_rated"
}

func (e *Event) ExpectedLastEventNumber() int64 {
	return e.expectedLastEventNumber
}

func (e *Event) SetExpectedLastEventNumber(number uint64) {
	e.expectedLastEventNumber = int64(number)
}
