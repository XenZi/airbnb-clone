package stream

import (
	"context"
	"example/metrics_events"
	user_joined "example/metrics_events/user_joined"
	user_left "example/metrics_events/user_left"
	user_rated "example/metrics_events/user_rated"
	user_reserved "example/metrics_events/user_reserved"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"log"
)

type ESDBEventStream struct {
	client *esdb.Client
	group  string
	sub    *esdb.PersistentSubscription
}

func NewESDBEventStream(client *esdb.Client, group string) (EventStream, error) {
	opts := esdb.PersistentAllSubscriptionOptions{
		From: esdb.Start{},
	}
	err := client.CreatePersistentSubscriptionAll(context.Background(), group, opts)
	if err != nil {
		// persistent subscription group already exists
		log.Println(err)
	}
	eventStream := &ESDBEventStream{
		client: client,
		group:  group,
	}
	err = eventStream.subscribe()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return eventStream, nil
}

func (s *ESDBEventStream) Process(processFn func(metrics_events.Event) error) {
	for {
		e := s.sub.Recv()

		if e.EventAppeared != nil {
			streamEvent := e.EventAppeared.Event
			log.Println(streamEvent.EventType)
			var event metrics_events.Event
			switch streamEvent.EventType {
			case metrics_events.EventTypeUserJoined:
				event = user_joined.NewEmptyEvent()
			case metrics_events.EventTypeUserLeft:
				event = user_left.NewEmptyEvent()
			case metrics_events.EventTypeUserRated:
				event = user_rated.NewEmptyEvent()
			case metrics_events.EventTypeUserReserved:
				event = user_reserved.NewEmptyEvent()
			}
			if event == nil {
				log.Println("unknown event type")
				continue
			}
			event.SetNumber(streamEvent.EventNumber)
			err := event.FromJSON(streamEvent.Data)
			err = processFn(event)
			if err != nil {
				log.Println(err)
				s.sub.Nack(err.Error(), esdb.Nack_Retry, e.EventAppeared)
			} else {
				s.sub.Ack(e.EventAppeared)
			}

		}

		if e.SubscriptionDropped != nil {
			log.Println(e.SubscriptionDropped.Error)
			// retry subscription
			for err := s.subscribe(); err != nil; {
			}
		}
	}
}

func (s *ESDBEventStream) subscribe() error {
	opts := esdb.ConnectToPersistentSubscriptionOptions{}
	sub, err := s.client.ConnectToPersistentSubscriptionToAll(context.Background(), s.group, opts)
	if err != nil {
		return err
	}
	s.sub = sub
	return nil
}
