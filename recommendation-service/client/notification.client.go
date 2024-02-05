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

type NotificationClient struct {
	address        string
	client         *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
}

type Notification struct {
	Mail string `json:"mail"`
	Text string `json:"text"`
}

func NewNotificationClient(host, port string, client *http.Client, cb *gobreaker.CircuitBreaker) *NotificationClient {
	return &NotificationClient{
		address:        fmt.Sprintf("http://%s:%s", host, port),
		client:         client,
		circuitBreaker: cb,
	}
}

func (nc NotificationClient) SendNewNotificationToUser(ctx context.Context, message string, userEmail string, userID string) *errors.ErrorStruct {
	data := Notification{
		Mail: userEmail,
		Text: message,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return errors.NewError(err.Error(), 500)
	}
	requestBody := bytes.NewReader(jsonData)
	cbResp, err := nc.circuitBreaker.Execute(func() (interface{}, error) {
		log.Println(nc.address + "/" + userID)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, nc.address+"/"+userID, requestBody)
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
	log.Println("RESP2 ERR USER", baseResp.Error)
	return errors.NewError(baseResp.Error, baseResp.Status)

}
