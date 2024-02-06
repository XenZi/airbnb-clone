package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/sony/gobreaker"
	"go.opentelemetry.io/otel/trace"
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
	tracer         trace.Tracer
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

func NewReservationClient(host, port string, client *http.Client, circuitBreaker *gobreaker.CircuitBreaker, tracer trace.Tracer) *ReservationClient {
	return &ReservationClient{
		address:        fmt.Sprintf("http://%s:%s", host, port),
		client:         client,
		circuitBreaker: circuitBreaker,
		tracer:         tracer,
	}
}

func (rc ReservationClient) getReservationsForUser(ctx context.Context, id, role string) ([]Reservation, *errors.ErrorStruct) {
	ctx, span := rc.tracer.Start(ctx, "ReservationClient.GetReservationsForUser")
	defer span.End()
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
		return nil, errors.NewError(err.Error(), 500)
	}
	resp := cbResp.(*http.Response)
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		baseResp := ResponseData{}
		err := json.NewDecoder(resp.Body).Decode(&baseResp)
		if err != nil {
			return nil, errors.NewError(err.Error(), 500)
		}
		list := baseResp.Data
		return list, nil
	}
	baseResp := domain.BaseErrorHttpResponse{}
	erro := json.NewDecoder(resp.Body).Decode(&baseResp)
	if erro != nil {
		return nil, errors.NewError(erro.Error(), 500)
	}
	return nil, errors.NewError(baseResp.Error, baseResp.Status)
}

func (rc ReservationClient) getCancelRate(ctx context.Context, id string) (*float32, *errors.ErrorStruct) {
	ctx, span := rc.tracer.Start(ctx, "ReservationClient.GetCancelRateForUser")
	defer span.End()
	cbResp, err := rc.circuitBreaker.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, rc.address+"/percentage-cancelation/"+id, http.NoBody)
		if err != nil {
			return nil, err
		}
		return rc.client.Do(req)
	})
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	resp := cbResp.(*http.Response)
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		baseResp := domain.BaseHttpResponse{}
		err := json.NewDecoder(resp.Body).Decode(&baseResp)
		if err != nil {
			return nil, errors.NewError(err.Error(), 500)
		}
		rate := baseResp.Data.(float32)
		return &rate, nil
	}
	baseResp := domain.BaseErrorHttpResponse{}
	erro := json.NewDecoder(resp.Body).Decode(&baseResp)
	if erro != nil {
		return nil, errors.NewError(erro.Error(), 500)
	}
	return nil, errors.NewError(baseResp.Error, baseResp.Status)
}

func (rc ReservationClient) UserDeleteAllowed(ctx context.Context, id, role string) *errors.ErrorStruct {
	ctx, span := rc.tracer.Start(ctx, "ReservationClient.UserDeleteAllowed")
	defer span.End()
	list, err := rc.getReservationsForUser(ctx, id, role)
	if err != nil {
		return err
	}
	if len(list) != 0 {
		return errors.NewError(fmt.Sprintf("User by id: %v has open reservations", id), 401)
	}
	return nil
}

// requirements lowered for easier examples
func (rc ReservationClient) CheckDistinguished(ctx context.Context, id string) (bool, *errors.ErrorStruct) {
	ctx, span := rc.tracer.Start(ctx, "ReservationClient.CheckDistinguished")
	defer span.End()
	reqCounter := 0
	list, err := rc.getReservationsForUser(ctx, id, "Host")
	if err != nil {
		return false, err
	}
	numOfReservations, err := checkNumberOfPastReservations(list)
	if err != nil {
		return false, err
	}
	if numOfReservations >= 3 {
		reqCounter += 1
	}
	durOfReservations := checkPastReservationDuration(list)
	if durOfReservations > 10 {
		reqCounter += 1
	}
	cancelRate, erro := rc.getCancelRate(ctx, id)
	var margin float32 = 0.33
	if erro != nil {
		return false, erro
	}
	if *cancelRate <= margin {
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
