package events

import (
	"example/metrics_events"
	user_joined "example/metrics_events/user_joined"
	user_left "example/metrics_events/user_left"
	user_rated "example/metrics_events/user_rated"
	user_reserved "example/metrics_events/user_reserved"
	"log"
	"metrics_query/domain"
	"time"
)

type EventHandler struct {
	store domain.AccommodationStore
}

func NewEventHandler(store domain.AccommodationStore) EventHandler {
	return EventHandler{
		store: store,
	}
}

func (h EventHandler) Handle(event metrics_events.Event) error {
	switch e := event.(type) {
	case *user_joined.Event:
		accommodation, err := h.store.Read(e.AccommodationID)
		if err != nil {
			accommodation = domain.Accommodation{
				Id:                                 e.AccommodationID,
				OnScreenTime:                       0,
				NumberOfVisits:                     0,
				NotClosedEventTimeStamps:           make(map[string]string),
				LastAppliedUserJoinedEventNumber:   -1,
				LastAppliedUserLeftEventNumber:     -1,
				NumberOfReservations:               0,
				LastAppliedUserReservedEventNumber: -1,
				NumberOfRatings:                    0,
				LastAppliedUserRatedEventNumber:    -1,
			}
			//accommodation.NotClosedEventTimeStamps[e.CustomUUID] = e.JoinedAt
			//accommodation.NumberOfVisits += 1
			//accommodation.LastAppliedUserJoinedEventNumber = int64(e.Number())
			err := h.store.Create(accommodation)
			if err != nil {
				return err
			}
		}
		accommodation.NotClosedEventTimeStamps[e.CustomUUID] = e.JoinedAt
		log.Println("UUID JE: ", e.CustomUUID)
		log.Println(e.JoinedAt)
		log.Println(accommodation.NotClosedEventTimeStamps[e.CustomUUID])
		accommodation.NumberOfVisits += 1
		accommodation.LastAppliedUserJoinedEventNumber = int64(e.Number())
		log.Println(accommodation)
		err = h.store.Update(accommodation)
		if err != nil {
			log.Println(err)
		}
	case *user_left.Event:
		accommodation, err := h.store.Read(e.AccommodationID)
		if err != nil {
			return err
		}
		joinedTimeString := accommodation.NotClosedEventTimeStamps[e.CustomUUID]
		leftTimeString := e.LeftAt
		log.Println("leftString je :", e.LeftAt)
		layout := "2006-01-02 15:04"
		joinedTime, err := time.Parse(layout, joinedTimeString)
		if err != nil {
			return err
		}
		leftTime, err := time.Parse(layout, leftTimeString)
		log.Println(joinedTime)
		log.Println(leftTime)
		duration := leftTime.Sub(joinedTime).Minutes()
		log.Println(duration)
		accommodation.OnScreenTime += duration
		accommodation.LastAppliedUserLeftEventNumber = int64(e.Number())
		delete(accommodation.NotClosedEventTimeStamps, e.CustomUUID)
		err = h.store.Update(accommodation)
		if err != nil {
			log.Println(err)
		}
	case *user_rated.Event:
		accommodation, err := h.store.Read(e.AccommodationID)
		if err != nil {
			return err
		}
		accommodation.NumberOfRatings += 1
		accommodation.LastAppliedUserRatedEventNumber = int64(e.Number())
		err = h.store.Update(accommodation)
		if err != nil {
			log.Println(err)
		}
	case *user_reserved.Event:
		accommodation, err := h.store.Read(e.AccommodationID)
		if err != nil {
			return err
		}
		accommodation.NumberOfReservations += 1
		accommodation.LastAppliedUserReservedEventNumber = int64(e.Number())
		err = h.store.Update(accommodation)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}
