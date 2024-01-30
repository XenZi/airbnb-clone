package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/XenZi/airbnb-clone/accommodations-service/domain"
	"github.com/XenZi/airbnb-clone/accommodations-service/errors"
	"log"
	"net/http"

	"github.com/sony/gobreaker"
)

type ReservationsClient struct {
	address        string
	client         *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
}

type SendCreateAccommodationAvailiabilty struct {
	AccommodationID string                               `json:"accommodationId"`
	Location        string                               `json:"location"`
	DateRange       []domain.AvailableAccommodationDates `json:"dateRange"`
}

func NewReservationsClient(host, port string, client *http.Client, circuitBreaker *gobreaker.CircuitBreaker) *ReservationsClient {
	return &ReservationsClient{
		address:        fmt.Sprintf("http://%s:%s", host, port),
		client:         http.DefaultClient,
		circuitBreaker: circuitBreaker,
	}
}

func (rc ReservationsClient) SendCreatedReservationsAvailabilities(ctx context.Context, id string, accommodation domain.CreateAccommodation) *errors.ErrorStruct {
	log.Println(accommodation)
	reqData := SendCreateAccommodationAvailiabilty{
		AccommodationID: id,
		Location:        accommodation.Location,
		DateRange:       accommodation.AvailableAccommodationDates}
	jsonData, err := json.Marshal(reqData)
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
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		baseResp := domain.BaseHttpResponse{}
		err := json.NewDecoder(resp.Body).Decode(&baseResp)
		if err != nil {
			return errors.NewError(err.Error(), 500)
		}
		log.Println("Base resp valid", baseResp)
		return nil
	}
	baseResp := domain.BaseErrorHttpResponse{}
	err = json.NewDecoder(resp.Body).Decode(&baseResp)
	if err != nil {
		return errors.NewError(err.Error(), 500)
	}
	log.Println(baseResp)
	log.Println(baseResp.Error)
	return errors.NewError(baseResp.Error, baseResp.Status)

	return errors.NewError("Nothing to parse", 500)
}

func (rc ReservationsClient) CheckAvailabilityForAccommodations(ctx context.Context, accommodationIDs []string, dateRange []string) ([]string, *errors.ErrorStruct) {
	log.Println("Uslo u Check")
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

	cbResp, err := rc.circuitBreaker.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, rc.address+"/accommodations", requestBody)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")

		return rc.client.Do(req)

	})
	if err != nil {
		return nil, errors.NewError("Internal server error", http.StatusInternalServerError)
	}
	response := cbResp.(*http.Response)
	log.Println("odgovor je", response)
	if response.StatusCode == 201 {
		resp := domain.BaseHttpResponse{}
		errO := json.NewDecoder(response.Body).Decode(&resp)
		if errO != nil {
			return nil, errors.NewError("Error decoding json", 500)
		}
		if resp.Data == nil {
			var stringSlice []string
			// Use stringSlice or return it as needed
			return stringSlice, nil
		}
		log.Println(resp.Data)
		dataSlice, ok := resp.Data.([]interface{})
		if !ok {
			fmt.Println("Data is not a []interface{}")
			return nil, errors.NewError("Error slicing", 500)
		}

		stringSlice := make([]string, len(dataSlice))
		for i, item := range dataSlice {
			if str, isString := item.(string); isString {
				stringSlice[i] = str
			} else {
				fmt.Printf("Element at index %d is not a string\n", i)
			}
		}

		log.Println("zauzete akomodacije", stringSlice)
		return stringSlice, nil
	} else {
		resp := domain.BaseErrorHttpResponse{}
		errO := json.NewDecoder(response.Body).Decode(&resp)
		if errO != nil {
			return nil, errors.NewError("Error decoding json", 500)
		}
		log.Println(resp)

		return nil, errors.NewError(resp.Error, resp.Status)

	}

}
