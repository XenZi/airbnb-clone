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

type MailNotification struct {
	Text string `json:"text"`
	Email string `json:"email"`
}

type MailClient struct {
	address string
	client  *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
}

func NewMailClient(host, port string, client *http.Client, circuitBreaker *gobreaker.CircuitBreaker) *MailClient {
	return &MailClient{
		address:        fmt.Sprintf("http://%s:%s", host, port),
		client:         client,
		circuitBreaker: circuitBreaker,
	}
}

func(mc MailClient) SendMailNotification(ctx context.Context, notification domains.Notification, email string) *errors.ErrorStruct{
	mailNotification := MailNotification{
		Text: notification.Text,
		Email: email,
	}
	jsonData, err := json.Marshal(mailNotification)
	if err != nil {
		return errors.NewError(err.Error(), 500)
	}
	requestBody := bytes.NewReader(jsonData)
	cbResp, err := mc.circuitBreaker.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, mc.address+"/send-notification-information", requestBody)
		if err != nil {
			return nil, err
		}
		return mc.client.Do(req)
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