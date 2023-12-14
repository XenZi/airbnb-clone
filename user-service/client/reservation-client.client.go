package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sony/gobreaker"
	"log"
	"net/http"
)

type ReservationClient struct {
	address        string
	client         *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
}

func NewReservationClient(host, port string, client *http.Client, circuitBreaker *gobreaker.CircuitBreaker) *ReservationClient {
	return &ReservationClient{
		address:        fmt.Sprintf("http://%s:%s", host, port),
		client:         client,
		circuitBreaker: circuitBreaker,
	}
}

func (rc ReservationClient) CheckUserReservations(ctx context.Context, id string) (bool, error) {

	cbResp, err := rc.circuitBreaker.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, rc.address+"/user/"+id, http.NoBody)
		if err != nil {
			return nil, err
		}
		return rc.client.Do(req)
	})

	if err != nil {
		log.Println("Comm Error  ", err)
		return false, err
	}
	resp := cbResp.(*http.Response)
	anResp := domain.BaseErrorHttpResponse{}

	err = json.NewDecoder(resp.Body).Decode(&anResp)
	if err != nil {
		return false, err
	}
	log.Println(anResp)
	return true, nil
}

func (rc ReservationClient) CheckAccommodationReservations(ctx context.Context, id string) (bool, error) {

	cbResp, err := rc.circuitBreaker.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, rc.address+"/accommodations/reservations/"+id, http.NoBody)
		if err != nil {
			return nil, err
		}
		return rc.client.Do(req)
	})

	if err != nil {
		log.Println("Comm Error  ", err)
		return false, err
	}
	resp := cbResp.(*http.Response)
	anResp := domain.BaseErrorHttpResponse{}

	err = json.NewDecoder(resp.Body).Decode(&anResp)
	if err != nil {
		return false, err
	}
	log.Println(anResp)
	return true, nil
}
