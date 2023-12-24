package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

type NotificationClient struct {
	address        string
	client         *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
}
type ReservationNotification struct {
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt"`
	IsOpened  bool   `json:"isOpened"`
}

func NewNotificationClient(host, port string, client *http.Client, circuitBreaker *gobreaker.CircuitBreaker) *NotificationClient {
	return &NotificationClient{
		address:        fmt.Sprintf("http://%s:%s", host, port),
		client:         http.DefaultClient,
		circuitBreaker: circuitBreaker,
	}
}
func (nc NotificationClient) request(method, url string, payload interface{}) (*http.Response, error) {
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

	return nc.client.Do(httpReq)
}
func (nc NotificationClient) SendReservationCreatedNotification(ctx context.Context, userId, message string) {
	req := ReservationNotification{
		Text:      message,
		CreatedAt: time.Now().String(),
		IsOpened:  false,
	}
	reqURL := nc.address + userId
	res, err := nc.request(http.MethodPost, reqURL, req)
	if err != nil || res.StatusCode != 502 {
		log.Println(err)
		return
	}
	log.Println("Notification for reservation has be sent")

}
