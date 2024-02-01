package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/sony/gobreaker"
	"log"
	"net/http"
	"time"
	"user-service/domain"
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

func (rc ReservationClient) getReservationsForUser(ctx context.Context, id, role string) ([]Reservation, *errors.ErrorStruct) {
	path := "host"
	if role == "Guest" {
		path = "guest"
	}
	cbResp, err := rc.circuitBreaker.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, rc.address+"/user/"+path+"/"+id, http.NoBody)
		if err != nil {
			return nil, err
		}
		return rc.client.Do(req)
	})
	if err != nil {
		return nil, errors.NewError("internal error", 500)
	}
	resp := cbResp.(*http.Response)
	if resp.StatusCode != 200 {
		return nil, errors.NewError("communication error", resp.StatusCode)
	}
	var list []Reservation
	var responseData ResponseData
	erro := json.NewDecoder(resp.Body).Decode(&responseData)
	if erro != nil {
		return nil, errors.NewError("data error", 500)
	}
	return list, nil
}

func (rc ReservationClient) UserDeleteAllowed(ctx context.Context, id, role string) (bool, *errors.ErrorStruct) {
	list, err := rc.getReservationsForUser(ctx, id, role)
	if err != nil {
		return false, err
	}
	if len(list) != 0 {
		return false, errors.NewError("user has pend reservations", 401)
	}
	return true, nil
}

func (rc ReservationClient) getCanelRate(ctx context.Context, id string) (*bool, *errors.ErrorStruct) {
	cbResp, err := rc.circuitBreaker.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, rc.address+"/percentage-cancelation/"+"/"+id, http.NoBody)
		if err != nil {
			return nil, err
		}
		return rc.client.Do(req)
	})
	if err != nil {
		return nil, errors.NewError("internal error", 500)
	}
	resp := cbResp.(*http.Response)
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		baseResp := domain.BaseHttpResponse{}
		err := json.NewDecoder(resp.Body).Decode(&baseResp)
		if err != nil {
			return nil, errors.NewError(err.Error(), 500)
		}
		log.Println("Base resp valid", baseResp)
		rateData := baseResp.Data.(float32)
		rate := checkCancelationRate(rateData)
		return &rate, nil
	}
	baseResp := domain.BaseErrorHttpResponse{}
	err = json.NewDecoder(resp.Body).Decode(&baseResp)
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	log.Println(baseResp)
	log.Println(baseResp.Error)
	return nil, errors.NewError(baseResp.Error, baseResp.Status)
}

// TODO requirements lowered for easier examples
func (rc ReservationClient) CheckDistinguished(ctx context.Context, id string) (bool, *errors.ErrorStruct) {
	reqCounter := 0
	list, err := rc.getReservationsForUser(ctx, id, "Host")
	if err != nil {
		return false, err
	}
	cancelLT, err := rc.getCanelRate(ctx, id)
	numOfReservations, err := checkNumberOfPastReservations(list)
	if err != nil {
		return false, err
	}
	// Req 1
	if numOfReservations >= 3 {
		reqCounter += 1
	}
	durOfReservations := checkPastReservationDuration(list)
	// Req 2
	if durOfReservations > 10 {
		reqCounter += 1
	}
	// Req 3
	if *cancelLT {
		reqCounter += 1
	}
	if reqCounter > 2 {
		log.Println("Bad Host, 0 stars")
		return false, nil
	}
	log.Println("pretty good, 5/7")
	return true, nil
}

func checkNumberOfPastReservations(list []Reservation) (int, *errors.ErrorStruct) {
	counter := 0
	for _, res := range list {
		date, err := time.Parse("2006-01-02", res.EndDate)
		if err != nil {
			return 0, errors.NewError("cannot parse date format", 500)
		}
		log.Println(date)
		if date.Before(time.Now()) {
			counter += 1
		}
	}
	log.Println("COUNTER FOR OLD RESERVATIONS: ", counter)
	return counter, nil
}

func checkPastReservationDuration(list []Reservation) int {
	counter := 0
	for _, res := range list {
		counter += res.NumberOfDays
	}
	log.Println("DURATION FOR OLD RESERVATIONS: ", counter)

	return counter
}

// TODO cancelation rate with reservations
func checkCancelationRate(resp float32) bool {
	if resp <= 0.33 {
		return true
	}
	return false
}
