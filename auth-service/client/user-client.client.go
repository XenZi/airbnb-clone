package client

import (
	"auth-service/domains"
	"auth-service/errors"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

func (uc UserClient) SendCreatedUser(ctx context.Context, id string, user domains.RegisterUser) *errors.ErrorStruct {
	userForUserService := struct {
		ID        string `json:"id"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Residence string `json:"residence"`
		Role      string `json:"role"`
		Username  string `json:"username"`
		Age       int    `json:"age"`
	}{
		ID:        id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Residence: user.CurrentPlace,
		Username:  user.Username,
		Role:      user.Role,
		Age:       user.Age,
	}
	jsonData, err := json.Marshal(userForUserService)
	if err != nil {
		return errors.NewError(err.Error(), 500)
	}
	requestBody := bytes.NewReader(jsonData)
	cbResp, err := uc.circuitBreaker.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, uc.address+"/create", requestBody)
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

func (uc UserClient) SendUpdateCredentials(ctx context.Context,updatedUser domains.User) *errors.ErrorStruct {
	userForUserService := struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		Username  string `json:"username"`
	}{
		ID:        updatedUser.ID.Hex(),
		Email:     updatedUser.Email,
		Username:  updatedUser.Username,
	}
	jsonData, err := json.Marshal(userForUserService)
	if err != nil {
		return errors.NewError(err.Error(), 500)
	}
	requestBody := bytes.NewReader(jsonData)
	cbResp, err := uc.circuitBreaker.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, uc.address+"/creds/" + updatedUser.ID.Hex(), requestBody)
		if err != nil {
			log.Println(err)
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
			return errors.NewError(err.Error(), 500)
		}
		return nil
	}
	baseResp := domains.BaseErrorHttpResponse{}
	err = json.NewDecoder(resp.Body).Decode(&baseResp)
	if err != nil {
		return errors.NewError(err.Error(), 500)
	}
	return errors.NewError(baseResp.Error, baseResp.Status)

}