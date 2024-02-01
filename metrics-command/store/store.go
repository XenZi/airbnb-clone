package store

type EventStore interface {
	Store(stream string, eventType string, event []byte) error
	StoreAndExpectLastEventNumber(stream string, eventType string, event []byte, lastEventNumber uint64) error
}
