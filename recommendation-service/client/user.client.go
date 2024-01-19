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

type UserClient struct {
	address        string
	client         *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
}

type SendingUserRating struct {
	Rating string `json:"rating"`
}

func NewUserClient(host, port string, client *http.Client, cb *gobreaker.CircuitBreaker) *UserClient {
	return &UserClient{
		address:        fmt.Sprintf("http://%s:%s", host, port),
		client:         client,
		circuitBreaker: cb,
	}
}

func (uc UserClient) SendNewRatingForUser(ctx context.Context, newValue float64, userID string) *errors.ErrorStruct {
	data := SendingUserRating{
		Rating: fmt.Sprintf("%v", newValue),
	}
	log.Println(data)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return errors.NewError(err.Error(), 500)
	}
	requestBody := bytes.NewReader(jsonData)
	cbResp, err := uc.circuitBreaker.Execute(func() (interface{}, error) {
		log.Println(uc.address + "/rating/" + userID)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, uc.address+"/rating/"+userID, requestBody)
		if err != nil {
			return nil, err
		}
		return uc.client.Do(req)
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
