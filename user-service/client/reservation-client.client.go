package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/sony/gobreaker"
	"net/http"
	"user-service/errors"
)

type ReservationClient struct {
	address        string
	client         *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
}

type ResponseData struct {
	Data []Reservation `json:"data"`
}
type Reservation struct {
	Id                gocql.UUID `json:"id"`
	UserID            string     `json:"userId"`
	AccommodationID   string     `json:"accommodationId"`
	StartDate         string     `json:"startDate"`
	EndDate           string     `json:"endDate"`
	Username          string     `json:"username"`
	AccommodationName string     `json:"accommodationName"`
	Location          string     `json:"location"`
	Price             int        `json:"price"`
	NumberOfDays      int        `json:"numOfDays"`
	Continent         string     `json:"continent"`
	DateRange         []string   `json:"dateRange"`
	IsActive          bool       `json:"isActive"`
	Country           string     `json:"country"`
	HostID            string     `json:"hostId"`
}

func NewReservationClient(host, port string, client *http.Client, circuitBreaker *gobreaker.CircuitBreaker) *ReservationClient {
	return &ReservationClient{
		address:        fmt.Sprintf("http://%s:%s", host, port),
		client:         client,
		circuitBreaker: circuitBreaker,
	}
}

func (rc ReservationClient) GuestDeleteAllowed(ctx context.Context, id string) (bool, *errors.ErrorStruct) {

	cbResp, err := rc.circuitBreaker.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, rc.address+"/user/guest/"+id, http.NoBody)
		if err != nil {
			return nil, err
		}
		return rc.client.Do(req)
	})
	if err != nil {
		return false, errors.NewError("internal error", 500)
	}
	resp := cbResp.(*http.Response)
	if resp.StatusCode != 200 {
		return false, errors.NewError("communication error", resp.StatusCode)
	}
	var list []Reservation
	var responseData ResponseData
	erro := json.NewDecoder(resp.Body).Decode(&responseData)
	if erro != nil {
		return false, errors.NewError("data error", 500)
	}
	list = responseData.Data
	if len(list) != 0 {
		return false, errors.NewError("user has pend reservations", 401)
	}
	return true, nil
}
