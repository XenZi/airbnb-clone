package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reservation-service/domain"
	"reservation-service/errors"

	"time"

	"github.com/sony/gobreaker"
)

type MetricsClient struct {
	address        string
	client         *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
}
type ReservationMetrics struct {
	UserID          string `json:"userId"`
	AccommodationID string `json:"accommodationId"`
	ReservedAt      string `json:"reservedAt"`
}

func NewMetricsClient(host, port string, client *http.Client, circuitBreaker *gobreaker.CircuitBreaker) *MetricsClient {
	return &MetricsClient{
		address:        fmt.Sprintf("http://%s:%s", host, port),
		client:         http.DefaultClient,
		circuitBreaker: circuitBreaker,
	}
}
func (mc MetricsClient) SendReserved(ctx context.Context, userID, accommodationID string) *errors.ReservationError {
	log.Println(fmt.Sprintf("vrednosti: %v %v", userID, accommodationID))

	now := time.Now()
	formattedTime := now.Format("2006-01-02 15:04")
	fmt.Println("Formatted Time:", formattedTime)

	metrics := ReservationMetrics{
		UserID:          userID,
		AccommodationID: accommodationID,
		ReservedAt:      formattedTime,
	}
	jsonData, err := json.Marshal(metrics)
	if err != nil {
		return errors.NewReservationError(500, err.Error())
	}
	log.Println("NEKIDATA:", jsonData)
	requestBody := bytes.NewReader(jsonData)
	cbResp, err := mc.circuitBreaker.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, mc.address+"/reserved", requestBody)
		if err != nil {
			return nil, err
		}
		return mc.client.Do(req)
	})
	if err != nil {
		return errors.NewReservationError(500, err.Error())
	}
	resp := cbResp.(*http.Response)
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		baseResp := domain.BaseHttpResponse{}
		err := json.NewDecoder(resp.Body).Decode(&baseResp)
		if err != nil {
			return errors.NewReservationError(500, err.Error())

		}
		log.Println("Base resp valid", baseResp)
		return nil
	}
	baseResp := domain.BaseErrorHttpResponse{}
	err = json.NewDecoder(resp.Body).Decode(&baseResp)
	if err != nil {
		return errors.NewReservationError(500, err.Error())

	}
	log.Println(baseResp)
	log.Println(baseResp.Error)
	return errors.NewReservationError(500, err.Error())

}
