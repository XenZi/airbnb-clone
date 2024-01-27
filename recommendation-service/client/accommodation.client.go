package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"recommendation-service/domains"
	"recommendation-service/errors"

	"github.com/sony/gobreaker"
)

type AccommodationClient struct {
	address        string
	client         *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
}

type SendingAccommodationRating struct {
	ID     string  `json:"id"`
	Rating float64 `json:"rating"`
}

type SendingListAccommodations struct {
	Ids []string `json:"ids"`
}

func NewAccommodationClient(host, port string, client *http.Client, cb *gobreaker.CircuitBreaker) *AccommodationClient {
	return &AccommodationClient{
		address:        fmt.Sprintf("http://%s:%s", host, port),
		client:         client,
		circuitBreaker: cb,
	}
}

func (ac AccommodationClient) SendNewRatingForAccommodation(ctx context.Context, newValue float64, accommodationID string) *errors.ErrorStruct {
	log.Print(ac.address)
	data := SendingAccommodationRating{
		ID:     accommodationID,
		Rating: newValue,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return errors.NewError(err.Error(), 500)
	}
	requestBody := bytes.NewReader(jsonData)
	cbResp, err := ac.circuitBreaker.Execute(func() (interface{}, error) {
		log.Println(ac.address + "/rating/" + accommodationID)
		req, err := http.NewRequestWithContext(ctx, http.MethodPut, ac.address+"/rating/"+accommodationID, requestBody)
		if err != nil {
			return nil, err
		}
		return ac.client.Do(req)
	})
	if err != nil {
		return errors.NewError(err.Error(), 500)
	}
	resp := cbResp.(*http.Response)
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		baseResp := domains.BaseHttpResponse{}
		err := json.NewDecoder(resp.Body).Decode(&baseResp)
		if err != nil {
			return errors.NewError(err.Error(), 500)
		}
		log.Println("Base resp valid", baseResp)
		return nil
	}
	baseResp := domains.BaseErrorHttpResponse{}
	err = json.NewDecoder(resp.Body).Decode(&baseResp)
	if err != nil {
		return errors.NewError(err.Error(), 500)
	}
	log.Println(baseResp)
	log.Println(baseResp.Error)
	return errors.NewError(baseResp.Error, baseResp.Status)
}

func (ac AccommodationClient) GetAllRecommendedAccommodationData(ctx context.Context, recommended []domains.Recommendation) ([]domains.AccommodationDTO, *errors.ErrorStruct) {
	log.Print(ac.address)
	structForSending := SendingListAccommodations{
		Ids: []string{},
	}
	for _, recommendation := range recommended {
		structForSending.Ids = append(structForSending.Ids, recommendation.AccommodationID)
	}
	log.Println(structForSending.Ids)
	jsonData, err := json.Marshal(structForSending)
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	requestBody := bytes.NewReader(jsonData)
	cbResp, err := ac.circuitBreaker.Execute(func() (interface{}, error) {
		log.Println(ac.address + "/recommended")
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, ac.address+"/recommended", requestBody)
		if err != nil {
			return nil, err
		}
		return ac.client.Do(req)
	})
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	resp := cbResp.(*http.Response)
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		baseResp := domains.BaseHttpResponse{}
		err := json.NewDecoder(resp.Body).Decode(&baseResp)
		if err != nil {
			return nil, errors.NewError(err.Error(), 500)
		}
		log.Println("Base resp valid", baseResp.Data)
		accommodations := make([]domains.AccommodationDTO, 0)
		for _, item := range baseResp.Data.([]interface{}) {
			itemJSON, err := json.Marshal(item)
			if err != nil {
				log.Println("Error marshaling item:", err)
				continue
			}

			var accommodation domains.AccommodationDTO
			err = json.Unmarshal(itemJSON, &accommodation)
			if err != nil {
				log.Println("Error unmarshaling item into AccommodationDTO:", err)
				continue
			}

			accommodations = append(accommodations, accommodation)
		}
		return accommodations, nil
	}
	baseResp := domains.BaseErrorHttpResponse{}
	err = json.NewDecoder(resp.Body).Decode(&baseResp)
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	log.Println(baseResp)
	log.Println(baseResp.Error)
	return nil, errors.NewError(baseResp.Error, baseResp.Status)
}
