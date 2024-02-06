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

type CommandQueryClient struct {
	address        string
	client         *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
}

type Rated struct {
	userID          string `json:"userID"`
	accommodationID string `json:"accommodationID"`
	ratedAt         string `json:"ratedAt"`
}

func NewCommandQueryClient(host, port string, client *http.Client, cb *gobreaker.CircuitBreaker) *CommandQueryClient {
	return &CommandQueryClient{
		address:        fmt.Sprintf("http://%s:%s", host, port),
		client:         client,
		circuitBreaker: cb,
	}
}

func (cq CommandQueryClient) SendNewRatingForAccommodation(ctx context.Context, userID, accommodationID, ratedAt string) *errors.ErrorStruct {
	data := Rated{
		userID:          userID,
		accommodationID: accommodationID,
		ratedAt:         ratedAt,
	}
	log.Println(data)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return errors.NewError(err.Error(), 500)
	}
	requestBody := bytes.NewReader(jsonData)
	cbResp, err := cq.circuitBreaker.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, cq.address+"/rated", requestBody)
		if err != nil {
			return nil, err
		}
		return cq.client.Do(req)
	})
	if err != nil {
		return errors.NewError(err.Error(), 500)
	}
	resp := cbResp.(*http.Response)
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		baseResp := domains.BaseHttpResponse{}
		err := json.NewDecoder(resp.Body).Decode(&baseResp)
		if err != nil {
			log.Println("AAAA", err.Error())
			return errors.NewError(err.Error(), 500)
		}
		log.Println("Base resp valid", baseResp)
		return nil
	}
	baseResp := domains.BaseErrorHttpResponse{}
	err = json.NewDecoder(resp.Body).Decode(&baseResp)
	if err != nil {
		log.Println("BBBB", err.Error())
		return errors.NewError(err.Error(), 500)
	}
	log.Println("RESP ERR USER", baseResp)
	log.Println("RESP2 ERR USER", baseResp.Error)
	return errors.NewError(baseResp.Error, baseResp.Status)

}
