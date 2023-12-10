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
	"time"
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
	log.Println(len(accommodation.AvailableAccommodationDates))
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

	}

	return errors.NewError("Nothing to parse", 500)
}

func (rc *ReservationsClient) GetAvailableAccommodations(ctx context.Context, startDateStr, endDateStr string, accommodationIDs []string) ([]string, *errors.ErrorStruct) {
	var availableAccommodationIDs []string
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, errors.NewError("Failed to parse start date", 500)
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return nil, errors.NewError("Failed to parse end date", 500)
	}

	requestData := struct {
		StartDate        string   `json:"startDate"`
		EndDate          string   `json:"endDate"`
		AccommodationIDs []string `json:"accommodationIds"`
	}{
		StartDate:        startDate.Format("2006-01-02"), // Format the dates as needed
		EndDate:          endDate.Format("2006-01-02"),
		AccommodationIDs: accommodationIDs,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, errors.NewError("Failed to parse request data", 500)
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
		log.Println("Error:", err)
		return nil, errors.NewError("Request failed", 500)
	}

	resp := cbResp.(*http.Response)
	defer resp.Body.Close()

	var response struct {
		AvailableAccommodations []string `json:"availableAccommodations"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, errors.NewError("Failed to parse response", 500)
	}

	availableAccommodationIDs = response.AvailableAccommodations

	return availableAccommodationIDs, nil
}
