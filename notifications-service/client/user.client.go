package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"notifications-service/domains"
	"notifications-service/errors"

	"github.com/sony/gobreaker"
)

type UserClient struct {
	address        string
	client         *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
}

func NewUserClient(host, port string, client *http.Client, circuitBreaker *gobreaker.CircuitBreaker) *UserClient {
	return &UserClient{
		address:        fmt.Sprintf("http://%s:%s", host, port),
		client:         client,
		circuitBreaker: circuitBreaker,
	}
}

func (uc UserClient) GetAllInformationsByUserID(ctx context.Context, id string) (*domains.User, *errors.ErrorStruct) {
	jsonData, err := json.Marshal(id)
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}
	requestBody := bytes.NewReader(jsonData)
	cbResp, err := uc.circuitBreaker.Execute(func() (interface{}, error) {
		log.Println(uc.address+"/"+id)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, uc.address+"/"+id, requestBody)
		if err != nil {
			log.Println("H ", err)
			return nil, err
		}
		return uc.client.Do(req)
	})
	if err != nil {
		log.Println("G ", err)

		return nil, errors.NewError(err.Error(), 500)
	}
	resp := cbResp.(*http.Response)
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		baseResp := domains.BaseHttpResponse{}
		err := json.NewDecoder(resp.Body).Decode(&baseResp)
		if err != nil {
			log.Println("A ", err)
			return nil, errors.NewError(err.Error(), 500)
		}
		log.Println("Base resp valid", baseResp)
	
		// Marshal the Data field back into JSON
		jsonData, err := json.Marshal(baseResp.Data)
		if err != nil {
			log.Println("B ", err)
			return nil, errors.NewError(err.Error(), 500)
		}
	
		// Unmarshal the JSON data into the domains.User struct
		var foundUser domains.User
		err = json.Unmarshal(jsonData, &foundUser)
		if err != nil {
			log.Println("C ", err)
			return nil, errors.NewError(err.Error(), 500)
		}
	
		return &foundUser, nil
	}
	baseResp := domains.BaseErrorHttpResponse{}
	err = json.NewDecoder(resp.Body).Decode(&baseResp)
	if err != nil {
		log.Println("S ", err)

		return nil, errors.NewError(err.Error(), 500)
	}
	log.Println(baseResp)
	log.Println(baseResp.Error)
	return nil, errors.NewError(baseResp.Error, baseResp.Status)

}