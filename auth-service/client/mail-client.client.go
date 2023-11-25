package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type MailClientInterface interface {
	SendAccountConfirmationEmail(email, token string)
	SendRequestResetPassword(email, token string)
}

type BaseClientSendStruct struct {
	Email string
	Token string
}

type MailClient struct {
	address string
	client  *http.Client
}

func NewMailClient(host, port string, client *http.Client) MailClientInterface {
	return &MailClient{
		address: fmt.Sprintf("http://%s:%s", host, port),
		client:  http.DefaultClient,
	}
}

func (mc MailClient) request(method, url string, payload interface{}) (*http.Response, error) {
	var bodyReader *bytes.Reader

	if payload != nil {
		reqBytes, err := json.Marshal(payload)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		bodyReader = bytes.NewReader(reqBytes)
	}

	httpReq, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return mc.client.Do(httpReq)
}

func (mc MailClient) SendAccountConfirmationEmail(email, token string) {
	req := BaseClientSendStruct{
		Email: email,
		Token: token,
	}
	requestURL := mc.address + "/confirm-new-account"
	res, err := mc.request(http.MethodPost, requestURL, req)
	if err != nil || res.StatusCode != 502 {
		log.Println(err)
		return
	}
	log.Println("Mail has been sent")
}

func (mc MailClient) SendRequestResetPassword(email, token string) {
	req := BaseClientSendStruct{
		Email: email,
		Token: token,
	}
	requestURL := mc.address + "/request-reset-password"
	_, err := mc.request(http.MethodPost, requestURL, req)
	log.Println(err)
	log.Println("Mail has been sent")
}