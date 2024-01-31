package handlers

import (
	"context"
	"errors"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/gorilla/mux"
	"io"
	"metrics_query/domain"
	"metrics_query/utils"
	"time"

	user_joined "example/metrics_events/user_joined"
	user_left "example/metrics_events/user_left"
	user_rated "example/metrics_events/user_rated"
	user_reserved "example/metrics_events/user_reserved"
	"net/http"
)

type AccommodationHandler struct {
	store  domain.AccommodationStore
	client *esdb.Client
}

func NewAccommodationHandler(store domain.AccommodationStore, client *esdb.Client) AccommodationHandler {
	return AccommodationHandler{
		store:  store,
		client: client,
	}
}

func (h AccommodationHandler) Get(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		utils.WriteErrorResp("bad request", 500, "metrics/get/{id}", rw)
		return
	}
	accommodation, err := h.store.Read(id)
	if err != nil {
		utils.WriteErrorResp(err.Error(), 404, "not found", rw)
		return
	}
	utils.WriteResp(accommodation, 200, rw)
	return
}

func (h AccommodationHandler) GetAll(rw http.ResponseWriter, r *http.Request) {
	accommodations := h.store.ReadAll()
	utils.WriteResp(accommodations, 200, rw)
}

func (h AccommodationHandler) Det(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		utils.WriteErrorResp("bad request", 500, "metrics/get/{id}", rw)
		return
	}
	accommodation := domain.Accommodation{
		Id:                                 id,
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
	var start uint64
	start = 0
	increment := 0
	stepMap := map[int]string{
		0: "user_joined",
		1: "user_left",
		2: "user_rated",
		3: "user_reserved",
	}
	for {
		val, _ := h.walkStream(start, stepMap[increment], &accommodation)
		if val {
			start += 100
			val, _ = h.walkStream(start, stepMap[increment], &accommodation)
		} else {
			start = 0
			increment += 1
			if increment == 5 {
				break
			}
			val, _ = h.walkStream(start, stepMap[increment], &accommodation)

		}
	}
	utils.WriteResp(accommodation, 200, rw)

}

func (h AccommodationHandler) walkStream(start uint64, streamName string, accommodation *domain.Accommodation) (bool, error) {
	stream, err := h.client.ReadStream(context.Background(), streamName, esdb.ReadStreamOptions{
		Direction: esdb.Backwards,
		From:      esdb.Revision(start),
	}, 100)
	if err != nil {
		return false, err
	}
	for {
		event, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return true, err
		}
		if err != nil {
			return false, err
		}
		switch streamName {
		case "user_joined":
			e := user_joined.NewEmptyEvent()
			err := e.FromJSON(event.Event.Data)
			if err != nil {
				return false, err
			}
			if !e.JoinedAt.After(time.Now().Add(-24 * time.Hour)) {
				return false, errors.New("out of bounds")
			}
			accommodation.NotClosedEventTimeStamps[e.CustomUUID] = e.JoinedAt
			accommodation.NumberOfVisits += 1
			accommodation.LastAppliedUserJoinedEventNumber = int64(e.Number())
		case "user_left":
			e := user_left.NewEmptyEvent()
			err := e.FromJSON(event.Event.Data)
			if err != nil {
				return false, errors.New("out of bounds")
			}
			if !e.LeftAt.After(time.Now().Add(-24 * time.Hour)) {
				return false, errors.New("out of bounds")
			}
			joinedTimeString := accommodation.NotClosedEventTimeStamps[e.CustomUUID]
			leftTimeString := e.LeftAt
			layout := "2006-01-02 15:04"
			joinedTime, err := time.Parse(layout, joinedTimeString)
			if err != nil {
				return false, err
			}
			leftTime, err := time.Parse(layout, leftTimeString)
			if err != nil {
				return false, err
			}
			duration := leftTime.Sub(joinedTime).Minutes()
			accommodation.OnScreenTime += duration
			accommodation.LastAppliedUserLeftEventNumber = int64(e.Number())
			delete(accommodation.NotClosedEventTimeStamps, e.CustomUUID)
		case "user_rated":
			e := user_rated.NewEmptyEvent()
			err := e.FromJSON(event.Event.Data)
			if err != nil {
				return false, err
			}
			if !e.RatedAt.After(time.Now().Add(-24 * time.Hour)) {
				return false, errors.New("out of bounds")
			}
			accommodation.NumberOfRatings += 1
			accommodation.LastAppliedUserRatedEventNumber = int64(e.Number())
		case "user_reserved":
			e := user_reserved.NewEmptyEvent()
			err := e.FromJSON(event.Event.Data)
			if err != nil {
				return false, err
			}
			if !e.ReservedAt.After(time.Now().Add(-24 * time.Hour)) {
				return false, errors.New("out of bounds")
			}
			accommodation.NumberOfReservations += 1
			accommodation.LastAppliedUserReservedEventNumber = int64(e.Number())
		}
	}
}
