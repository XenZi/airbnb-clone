package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sony/gobreaker"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"user-service/domain"
	"user-service/errors"
)

type AccClient struct {
	address        string
	client         *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
	tracer         trace.Tracer
}

func NewAccClient(host, port string, client *http.Client, circuitBreaker *gobreaker.CircuitBreaker, tracer trace.Tracer) *AccClient {
	return &AccClient{
		address:        fmt.Sprintf("http://%s:%s", host, port),
		client:         client,
		circuitBreaker: circuitBreaker,
		tracer:         tracer,
	}
}

func (ac AccClient) DeleteUserAccommodations(ctx context.Context, id string) *errors.ErrorStruct {
	ctx, span := ac.tracer.Start(ctx, "AccommodationClient.DeleteUserAccommodations")
	defer span.End()
	cbResp, err := ac.circuitBreaker.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, ac.address+"/user/"+id, http.NoBody)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		return ac.client.Do(req)
	})
	if err != nil {
		return errors.NewError(err.Error(), 500)
	}
	resp := cbResp.(*http.Response)
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return nil
	}
	baseResp := domain.BaseErrorHttpResponse{}
	erro := json.NewDecoder(resp.Body).Decode(&baseResp)
	if erro != nil {
		return errors.NewError(erro.Error(), 500)
	}
	return errors.NewError(baseResp.Error, baseResp.Status)
}
