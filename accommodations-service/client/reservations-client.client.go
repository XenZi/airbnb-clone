package client

import (
	"accommodations-service/config"
	"accommodations-service/domain"
	"accommodations-service/errors"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sony/gobreaker"
)

type ReservationsClient struct {
	address        string
	client         *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
	logger         *config.Logger
}

type SendCreateAccommodationAvailiabilty struct {
	AccommodationID string                               `json:"accommodationId"`
	Location        string                               `json:"location"`
	DateRange       []domain.AvailableAccommodationDates `json:"dateRange"`
}

func NewReservationsClient(host, port string, client *http.Client, circuitBreaker *gobreaker.CircuitBreaker, logger *config.Logger) *ReservationsClient {
	return &ReservationsClient{
		address:        fmt.Sprintf("http://%s:%s", host, port),
		client:         http.DefaultClient,
		circuitBreaker: circuitBreaker,
		logger:         logger,
	}
}

func (rc ReservationsClient) SendCreatedReservationsAvailabilities(ctx context.Context, id string, accommodation domain.CreateAccommodation) *errors.ErrorStruct {
	log.Println(accommodation)
	reqData := SendCreateAccommodationAvailiabilty{
		AccommodationID: id,
		Location:        accommodation.Location,
		DateRange:       accommodation.AvailableAccommodationDates,
	}
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		rc.logger.LogError("accommodations-client", fmt.Sprintf("Unable to parse json data"))
		rc.logger.LogError("accommodation-client", fmt.Sprintf("Error:"+err.Error()))
		return errors.NewError("Nothing to parse", 500)
	}

	requestBody := bytes.NewReader(jsonData)

	cbResp, err := rc.circuitBreaker.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, rc.address+"/availability", requestBody)
		if err != nil {
			rc.logger.LogError("accommodations-client", fmt.Sprintf("Unable to send request to reservations service"))
			rc.logger.LogError("accommodation-client", fmt.Sprintf("Error:"+err.Error()))
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
			rc.logger.LogError("accommodations-client", fmt.Sprintf("Unable to decode json data"))
			rc.logger.LogError("accommodation-client", fmt.Sprintf("Error:"+err.Error()))
			return errors.NewError(err.Error(), 500)
		}
		log.Println("Base resp valid", baseResp)
		return nil
	}
	baseResp := domain.BaseErrorHttpResponse{}
	err = json.NewDecoder(resp.Body).Decode(&baseResp)
	if err != nil {
		rc.logger.LogError("accommodations-client", fmt.Sprintf("Unable to decode json data"))
		rc.logger.LogError("accommodation-client", fmt.Sprintf("Error:"+err.Error()))
		return errors.NewError(err.Error(), 500)
	}
	log.Println(baseResp)
	log.Println(baseResp.Error)
	rc.logger.LogInfo("accommodation-client", fmt.Sprintf("Successfully sent availabilities"))
	return errors.NewError(baseResp.Error, baseResp.Status)

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
		rc.logger.LogError("accommodations-client", fmt.Sprintf("Unable to marshal json data"))
		rc.logger.LogError("accommodation-client", fmt.Sprintf("Error:"+err.Error()))
		return nil, errors.NewError("Failed to marshal JSON data", http.StatusInternalServerError)
	}

	requestBody := bytes.NewReader(jsonData)

	cbResp, err := rc.circuitBreaker.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, rc.address+"/accommodations", requestBody)
		if err != nil {
			rc.logger.LogError("accommodations-client", fmt.Sprintf("Unable to send request"))
			rc.logger.LogError("accommodation-client", fmt.Sprintf("Error:"+err.Error()))
			log.Println(err)
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		rc.logger.LogInfo("accommodation-client", fmt.Sprintf("Successfully sent request to reservations server"))
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
			rc.logger.LogError("accommodations-client", fmt.Sprintf("Unable to decode response"))
			rc.logger.LogError("accommodation-client", fmt.Sprintf("Error:"+errO.Error()))
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
			rc.logger.LogError("accommodations-client", fmt.Sprintf("Error in maping data"))
			rc.logger.LogError("accommodation-client", fmt.Sprintf("Error:"+err.Error()))
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
			rc.logger.LogError("accommodations-client", fmt.Sprintf("Unable to decode json"))
			rc.logger.LogError("accommodation-client", fmt.Sprintf("Error:"+err.Error()))
			return nil, errors.NewError("Error decoding json", 500)
		}
		log.Println(resp)
		rc.logger.LogInfo("accommodation-client", fmt.Sprintf("Successfully decoded and collected data sent from reservations server"))
		return nil, errors.NewError(resp.Error, resp.Status)

	}

}

func (rc ReservationsClient) GetAccommodationsBelowPrice(ctx context.Context, maxPrice int) ([]string, *errors.ErrorStruct) {
	log.Println("Entering GetAccommodationsBelowPrice")

	// Build the request URL with the max price
	url := fmt.Sprintf("%s/price/myPrice/%d", rc.address, maxPrice)

	cbResp, err := rc.circuitBreaker.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			rc.logger.LogError("accommodations-client", fmt.Sprintf("Unable to send request to reservations server"))
			rc.logger.LogError("accommodation-client", fmt.Sprintf("Error:"+err.Error()))
			log.Println(err)
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		rc.logger.LogInfo("accommodation-client", fmt.Sprintf("Successfully sent request to reservations server"))
		return rc.client.Do(req)
	})

	if err != nil {
		rc.logger.LogError("accommodations-client", fmt.Sprintf("Internal server error"))
		rc.logger.LogError("accommodation-client", fmt.Sprintf("Error:"+err.Error()))
		return nil, errors.NewError("Internal server error", http.StatusInternalServerError)
	}

	response := cbResp.(*http.Response)
	log.Println("Response received:", response)

	if response.StatusCode == http.StatusOK {
		resp := domain.BaseHttpResponse{}
		err := json.NewDecoder(response.Body).Decode(&resp)
		if err != nil {
			rc.logger.LogError("accommodations-client", fmt.Sprintf("Error decoding response body"))
			rc.logger.LogError("accommodation-client", fmt.Sprintf("Error:"+err.Error()))
			return nil, errors.NewError("Error decoding JSON", http.StatusInternalServerError)
		}

		if resp.Data == nil {
			var stringSlice []string
			// Use stringSlice or return it as needed
			return stringSlice, nil
		}

		log.Println(resp.Data)
		dataSlice, ok := resp.Data.([]interface{})
		if !ok {
			rc.logger.LogError("accommodations-client", fmt.Sprintf("Internal server error"))
			rc.logger.LogError("accommodation-client", fmt.Sprintf("Error:"+err.Error()))
			fmt.Println("Data is not a []interface{}")
			return nil, errors.NewError("Error slicing", http.StatusInternalServerError)
		}

		stringSlice := make([]string, len(dataSlice))
		for i, item := range dataSlice {
			if str, isString := item.(string); isString {
				stringSlice[i] = str
			} else {
				fmt.Printf("Element at index %d is not a string\n", i)
			}
		}

		log.Println("Accommodations below price:", stringSlice)
		rc.logger.LogInfo("accommodation-client", fmt.Sprintf("Successfully decoded and collected data sent from reservations server"))
		return stringSlice, nil
	} else {
		resp := domain.BaseErrorHttpResponse{}
		err := json.NewDecoder(response.Body).Decode(&resp)
		if err != nil {
			rc.logger.LogError("accommodations-client", fmt.Sprintf("Error decoding response body"))
			rc.logger.LogError("accommodation-client", fmt.Sprintf("Error:"+err.Error()))
			return nil, errors.NewError("Error decoding JSON", http.StatusInternalServerError)
		}

		log.Println(resp)
		return nil, errors.NewError(resp.Error, resp.Status)
	}
}
