package store

import (
	"context"
	"errors"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/gofrs/uuid"
)

var (
	ErrEmptyStream    = errors.New("no events in the stream")
	ErrStreamNOtFound = errors.New("stream not found")
)

type ESDBStore struct {
	client *esdb.Client
}

func NewESDBStore(client *esdb.Client) EventStore {
	return ESDBStore{
		client: client,
	}
}

func (e ESDBStore) Store(stream string, eventType string, event []byte) error {
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	eventData := esdb.EventData{
		EventID:     id,
		EventType:   eventType,
		Data:        event,
		ContentType: esdb.JsonContentType,
	}
	opts := esdb.AppendToStreamOptions{}
	_, err = e.client.AppendToStream(context.Background(), stream, opts, eventData)
	return err
}

func (e ESDBStore) StoreAndExpectLastEventNumber(stream string, eventType string, event []byte, lastEventNumber uint64) error {
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	eventData := esdb.EventData{
		EventID:     id,
		EventType:   eventType,
		Data:        event,
		ContentType: esdb.JsonContentType,
	}
	opts := esdb.AppendToStreamOptions{
		ExpectedRevision: esdb.Revision(lastEventNumber),
	}
	_, err = e.client.AppendToStream(context.Background(), stream, opts, eventData)
	return err
}
