package metrics_events

const (
	EventTypeUserJoined   = "UserJoined"
	EventTypeUserLeft     = "UserLeft"
	EventTypeUserRated    = "UserRated"
	EventTypeUserReserved = "UserReserved"
)

type Event interface {
	Type() string
	ToJSON() ([]byte, error)
	FromJSON(jsonEvent []byte) error
	Number() uint64
	SetNumber(number uint64)
	Stream() string
	ExpectedLastEventNumber() int64
	SetExpectedLastEventNumber(number uint64)
}
