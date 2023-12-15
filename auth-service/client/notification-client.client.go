package client

import (
	"auth-service/domains"
	"auth-service/errors"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sony/gobreaker"
)

type NotificationClient struct {
	address string
	client  *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
}

func NewNotificationClient(host, port string, client *http.Client, circuitBreaker *gobreaker.CircuitBreaker) *NotificationClient {
	return &NotificationClient{
		address:        fmt.Sprintf("http://%s:%s", host, port),
		client:         client,
		circuitBreaker: circuitBreaker,
	}
}

func(nc NotificationClient) CreateNewUserStructNotification(ctx context.Context, id string) *errors.ErrorStruct {
	
	cbResp, err := nc.circuitBreaker.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, nc.address+"/create-new-user-notification/"+id, nil)
		if err != nil {
			return nil, err
		}
		return nc.client.Do(req)
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