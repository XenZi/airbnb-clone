package client

import (
	"accommodations-service/domain"
	"accommodations-service/errors"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sony/gobreaker"
	"log"
	"net/http"
)

type ReservationsClient struct {
	address        string
	client         *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
}

func NewReservationsClient(host, port string, client *http.Client, circuitBreaker *gobreaker.CircuitBreaker) *ReservationsClient {
	return &ReservationsClient{
		address:        fmt.Sprintf("http://%s:%s", host, port),
		client:         http.DefaultClient,
		circuitBreaker: circuitBreaker,
	}
}

func (rc ReservationsClient) SendCreatedReservationsAvailabilities(ctx context.Context, id string, accommodation domain.CreateAccommodation) *errors.ErrorStruct {
	for i := 0; i < len(accommodation.AvailableAccommodationDates); i++ {

		log.Println("OVO JE ID", accommodation.AvailableAccommodationDates[i].StartDate)
		availabilitiesForReservationService := struct {
			AccommodationID string `json:"accommodationId"`
			StartDate       string `json:"startDate"`
			EndDate         string `json:"endDate"`
			Location        string `json:"location"`
			Price           int    `json:"price"`
		}{
			AccommodationID: id,
			StartDate:       accommodation.AvailableAccommodationDates[i].StartDate,
			EndDate:         accommodation.AvailableAccommodationDates[i].EndDate,
			Location:        accommodation.Location,
			Price:           accommodation.AvailableAccommodationDates[i].Price,
		}

		jsonData, err := json.Marshal(availabilitiesForReservationService)
		if err != nil {
			return errors.NewError("Nothing to parse", 500)
		}
		requestBody := bytes.NewReader(jsonData)
		cbResp, err := rc.circuitBreaker.Execute(func() (interface{}, error) {
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, rc.address+"/availability", requestBody)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			return rc.client.Do(req)
		})

		if err != nil {
			log.Println("ERR FROM GGG ", err)
			return errors.NewError("Nothing to parse", 500)
		}
		resp := cbResp.(*http.Response)
		anResp := domain.BaseErrorHttpResponse{}

		err = json.NewDecoder(resp.Body).Decode(&anResp)
		if err != nil {
			return errors.NewError("Nothing to parse", 500)
		}
		log.Println(anResp)
		return nil
	}

	return errors.NewError("Nothing to parse", 500)
}
