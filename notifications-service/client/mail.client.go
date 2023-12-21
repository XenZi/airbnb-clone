package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"notifications-service/domains"
)

type MailNotification struct {
	Text string `json:"text"`
	Email string `json:"email"`
}

type MailClient struct {
	address string
	client  *http.Client
}

func NewMailClient(host, port string, client *http.Client) *MailClient {
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

func (mc MailClient) SendMailNotification(notification domains.Notification, email string) {
	req := MailNotification{
		Email: email,
		Text: notification.Text,
	}
	requestURL := mc.address + "/send-notification-information"
	res, err := mc.request(http.MethodPost, requestURL, req)
	if err != nil || res.StatusCode != 502 {
		log.Println(err)
		return
	}
	log.Println("Mail has been sent")
}