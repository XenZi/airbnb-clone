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
	daily := "daily"
	monthly := "monthly"
	switch e := event.(type) {
	case *user_joined.Event:
		eventDateStr, err := getTimeFromString(e.JoinedAt)
		if err != nil {
			return err
		}
		eventDay := getDayStart(*eventDateStr)
		accommodation, err := h.store.Read(e.AccommodationID, daily)
		if err != nil {
			accommodation = generateBlankReport(eventDay, e.AccommodationID)
			accommodation.NotClosedEventTimeStamps[e.CustomUUID] = e.JoinedAt
			accommodation.NumberOfVisits += 1
			accommodation.LastAppliedUserJoinedEventNumber = int64(e.Number())
			err := h.store.Create(*accommodation, daily)
			if err != nil {
				return err
			}
			err = h.store.Create(*accommodation, monthly)
			if err != nil {
				return err
			}
			return nil
		}
		monthlyAccommodation, err := h.store.Read(e.AccommodationID, monthly)
		if checkDay(accommodation.ReportingDate, eventDay) {
			accommodation.NotClosedEventTimeStamps[e.CustomUUID] = e.JoinedAt
			accommodation.NumberOfVisits += 1
			accommodation.LastAppliedUserJoinedEventNumber = int64(e.Number())
			err = h.store.Update(*accommodation, daily)
			if err != nil {
				log.Println(err)
			}
			err = h.store.Update(*accommodation, monthly)
			if err != nil {
				log.Println(err)
			}
		} else if checkMonth(monthlyAccommodation.ReportingDate, eventDay) {
			accommodation := generateBlankReport(eventDay, e.AccommodationID)
			accommodation.NotClosedEventTimeStamps[e.CustomUUID] = e.JoinedAt
			accommodation.NumberOfVisits += 1
			accommodation.LastAppliedUserJoinedEventNumber = int64(e.Number())
			err = h.store.Update(*accommodation, daily)
			if err != nil {
				return err
			}
			monthlyAccommodation.NotClosedEventTimeStamps[e.CustomUUID] = e.JoinedAt
			monthlyAccommodation.NumberOfVisits += 1
			monthlyAccommodation.LastAppliedUserJoinedEventNumber = int64(e.Number())
			err = h.store.Update(*monthlyAccommodation, monthly)
			if err != nil {
				return err
			}
		} else {
			accommodation := generateBlankReport(eventDay, e.AccommodationID)
			accommodation.NotClosedEventTimeStamps[e.CustomUUID] = e.JoinedAt
			accommodation.NumberOfVisits += 1
			accommodation.LastAppliedUserJoinedEventNumber = int64(e.Number())
			err = h.store.Update(*accommodation, daily)
			if err != nil {
				return err
			}
			err = h.store.Update(*accommodation, monthly)
			if err != nil {
				log.Println(err)
			}
		}

	case *user_left.Event:
		eventDateStr, err := getTimeFromString(e.LeftAt)
		if err != nil {
			return err
		}
		eventDay := getDayStart(*eventDateStr)
		accommodation, err := h.store.Read(e.AccommodationID, daily)
		if err != nil {
			return err
		}
		monthlyAccommodation, err := h.store.Read(e.AccommodationID, monthly)
		if err != nil {
			return err
		}
		joinedTime, err := getTimeFromString(accommodation.NotClosedEventTimeStamps[e.CustomUUID])
		if err != nil {
			return err
		}
		leftTime, err := getTimeFromString(e.LeftAt)
		if err != nil {
			return err
		}
		duration := leftTime.Sub(*joinedTime).Minutes()
		if checkDay(accommodation.ReportingDate, eventDay) {
			accommodation.OnScreenTime += duration
			accommodation.LastAppliedUserLeftEventNumber = int64(e.Number())
			delete(accommodation.NotClosedEventTimeStamps, e.CustomUUID)
			err = h.store.Update(*accommodation, daily)
			if err != nil {
				log.Println(err)
			}
			monthlyAccommodation.OnScreenTime += duration
			monthlyAccommodation.LastAppliedUserLeftEventNumber = int64(e.Number())
			delete(monthlyAccommodation.NotClosedEventTimeStamps, e.CustomUUID)
			err = h.store.Update(*accommodation, monthly)
			if err != nil {
				log.Println(err)
			}
		} else if checkMonth(monthlyAccommodation.ReportingDate, eventDay) {
			accommodation = generateBlankReport(eventDay, e.AccommodationID)
			accommodation.OnScreenTime += duration / 2
			accommodation.NumberOfVisits += 1
			accommodation.LastAppliedUserJoinedEventNumber = monthlyAccommodation.LastAppliedUserJoinedEventNumber
			accommodation.LastAppliedUserLeftEventNumber = int64(e.Number())
			err = h.store.Update(*accommodation, daily)
			if err != nil {
				log.Println(err)
			}
			monthlyAccommodation.OnScreenTime += duration
			monthlyAccommodation.LastAppliedUserLeftEventNumber = int64(e.Number())
			delete(monthlyAccommodation.NotClosedEventTimeStamps, e.CustomUUID)
			err = h.store.Update(*accommodation, monthly)
			if err != nil {
				log.Println(err)
			}
		} else {
			accommodation = generateBlankReport(eventDay, e.AccommodationID)
			accommodation.OnScreenTime += duration / 2
			accommodation.NumberOfVisits += 1
			accommodation.LastAppliedUserJoinedEventNumber = monthlyAccommodation.LastAppliedUserJoinedEventNumber
			accommodation.LastAppliedUserLeftEventNumber = int64(e.Number())
			err = h.store.Update(*accommodation, daily)
			if err != nil {
				log.Println(err)
			}
			err = h.store.Update(*accommodation, monthly)
			if err != nil {
				log.Println(err)
			}
		}

	case *user_rated.Event:
		eventDateStr, err := getTimeFromString(e.RatedAt)
		if err != nil {
			return err
		}
		eventDay := getDayStart(*eventDateStr)
		accommodation, err := h.store.Read(e.AccommodationID, daily)
		if err != nil {
			return err
		}
		monthlyAccommodation, err := h.store.Read(e.AccommodationID, monthly)
		if err != nil {
			return err
		}
		if checkDay(accommodation.ReportingDate, eventDay) {
			accommodation.NumberOfRatings += 1
			accommodation.LastAppliedUserRatedEventNumber = int64(e.Number())
			err = h.store.Update(*accommodation, daily)
			if err != nil {
				log.Println(err)
			}
			monthlyAccommodation.NumberOfRatings += 1
			monthlyAccommodation.LastAppliedUserRatedEventNumber = int64(e.Number())
			err = h.store.Update(*monthlyAccommodation, monthly)
			if err != nil {
				log.Println(err)
			}
		} else if checkMonth(monthlyAccommodation.ReportingDate, eventDay) {
			accommodation = generateBlankReport(eventDay, e.AccommodationID)
			accommodation.NumberOfRatings += 1
			accommodation.LastAppliedUserRatedEventNumber = int64(e.Number())
			err = h.store.Update(*accommodation, daily)
			if err != nil {
				log.Println(err)
			}
			accommodation.NumberOfRatings += 1
			accommodation.LastAppliedUserRatedEventNumber = int64(e.Number())
			err = h.store.Update(*accommodation, monthly)
			if err != nil {
				log.Println(err)
			}
		} else {
			accommodation = generateBlankReport(eventDay, e.AccommodationID)
			accommodation.NumberOfRatings += 1
			accommodation.LastAppliedUserRatedEventNumber = int64(e.Number())
			err = h.store.Update(*accommodation, daily)
			if err != nil {
				log.Println(err)
			}
			err = h.store.Update(*accommodation, monthly)
			if err != nil {
				log.Println(err)
			}
		}
	case *user_reserved.Event:
		eventDateStr, err := getTimeFromString(e.ReservedAt)
		if err != nil {
			return err
		}
		eventDay := getDayStart(*eventDateStr)
		accommodation, err := h.store.Read(e.AccommodationID, daily)
		if err != nil {
			return err
		}
		monthlyAccommodation, err := h.store.Read(e.AccommodationID, monthly)
		if err != nil {
			return err
		}
		if checkDay(accommodation.ReportingDate, eventDay) {
			accommodation.NumberOfReservations += 1
			accommodation.LastAppliedUserReservedEventNumber = int64(e.Number())
			err = h.store.Update(*accommodation, daily)
			if err != nil {
				log.Println(err)
			}
			monthlyAccommodation.NumberOfReservations += 1
			monthlyAccommodation.LastAppliedUserReservedEventNumber = int64(e.Number())
			err = h.store.Update(*monthlyAccommodation, monthly)
			if err != nil {
				log.Println(err)
			}
		} else if checkMonth(monthlyAccommodation.ReportingDate, eventDay) {
			accommodation = generateBlankReport(eventDay, e.AccommodationID)
			accommodation.NumberOfReservations += 1
			accommodation.LastAppliedUserReservedEventNumber = int64(e.Number())
			err = h.store.Update(*accommodation, daily)
			if err != nil {
				log.Println(err)
			}
			monthlyAccommodation.NumberOfReservations += 1
			monthlyAccommodation.LastAppliedUserReservedEventNumber = int64(e.Number())
			err = h.store.Update(*accommodation, monthly)
			if err != nil {
				log.Println(err)
			}
		} else {
			accommodation = generateBlankReport(eventDay, e.AccommodationID)
			accommodation.NumberOfReservations += 1
			accommodation.LastAppliedUserReservedEventNumber = int64(e.Number())
			err = h.store.Update(*accommodation, daily)
			if err != nil {
				log.Println(err)
			}
			err = h.store.Update(*accommodation, monthly)
			if err != nil {
				log.Println(err)
			}
		}
	}
	return nil
}

func getTimeFromString(str string) (*time.Time, error) {
	layout := "2006-01-02 15:04"
	tm, err := time.Parse(layout, str)
	if err != nil {
		return nil, err
	}
	return &tm, nil
}

func getDayStart(tm time.Time) time.Time {
	year, month, date := tm.Date()
	return time.Date(year, month, date, 0, 0, 0, 0, tm.Location())
}

func checkDay(old, new time.Time) bool {
	y1, m1, d1 := old.Date()
	y2, m2, d2 := new.Date()
	if y1 == y2 && m1 == m2 && d1 == d2 {
		return true
	}
	return false
}

func checkMonth(old, new time.Time) bool {
	y1, m1, _ := old.Date()
	y2, m2, _ := new.Date()
	if y1 == y2 && m1 == m2 {
		return true
	}
	return false
}

func generateBlankReport(tm time.Time, id string) *domain.Accommodation {
	accommodation := domain.Accommodation{
		Id:                                 id,
		ReportingDate:                      tm,
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
	return &accommodation
}
