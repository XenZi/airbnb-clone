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

func (rc ReservationsClient) CheckAvailabilityForAccommodations(ctx context.Context, accommodationIDs []string, dateRange []string) ([]string, *errors.ErrorStruct) {
	availabilityCheck := struct {
		AccommodationIDs []string `json:"accommodationIDs"`
		DateRange        []string `json:"dateRange"`
	}{
		AccommodationIDs: accommodationIDs,
		DateRange:        dateRange,
	}

	jsonData, err := json.Marshal(availabilityCheck)
	if err != nil {
		return nil, errors.NewError("Failed to marshal JSON data", http.StatusInternalServerError)
	}

	requestBody := bytes.NewReader(jsonData)

	var responseAccommodationIDs []string

	cbResp, err := rc.circuitBreaker.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, rc.address+"/availabilityFinder", requestBody)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := rc.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("non-success status code received: %d", resp.StatusCode)
		}

		var responseData struct {
			AccommodationIDs []string `json:"accommodationIDs"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
			return nil, err
		}

		responseAccommodationIDs = responseData.AccommodationIDs

		return resp, nil
	})

	if err != nil {
		log.Println("ERR FROM GGG ", err)
		return nil, errors.NewError("Failed to perform HTTP request", http.StatusInternalServerError)
	}

	resp := cbResp.(*http.Response)
	anResp := domain.BaseErrorHttpResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&anResp); err != nil {
		return nil, errors.NewError("Failed to parse response body", http.StatusInternalServerError)
	}
	log.Println(anResp)

	return responseAccommodationIDs, nil
}
