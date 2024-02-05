package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sony/gobreaker"
	"log"
	"net/http"
	"user-service/domain"
	"user-service/errors"
)

type AccClient struct {
	address        string
	client         *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
}

func NewAccClient(host, port string, client *http.Client, circuitBreaker *gobreaker.CircuitBreaker) *AccClient {
	return &AccClient{
		address:        fmt.Sprintf("http://%s:%s", host, port),
		client:         client,
		circuitBreaker: circuitBreaker,
	}
}

func (ac AccClient) DeleteUserAccommodations(ctx context.Context, id string) *errors.ErrorStruct {
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
