package client

import (
	"auth-service/domains"
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

func (uc UserClient) SendCreatedUser(ctx context.Context, id string, user domains.RegisterUser) error {
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
		return fmt.Errorf("error marshalling user data: %v", err)
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
		log.Println("ERR FROM GGG ", err)
		return err
	}
	resp := cbResp.(*http.Response)
	anResp := struct {
		ID        string `json:"id"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Residence string `json:"residence"`
		Role      string `json:"role"`
		Username  string `json:"username"`
		Age       int    `json:"age"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&anResp)
	if err != nil {
		return err
	}
	log.Println(anResp)
	return nil
}
