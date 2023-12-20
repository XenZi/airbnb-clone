package client

import (
	"context"
	"fmt"
	"github.com/sony/gobreaker"
	"log"
	"net/http"
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
			log.Println(err)
			return nil, err
		}
		return ac.client.Do(req)
	})
	if err != nil {
		return errors.NewError("internal error", 500)
	}
	resp := cbResp.(*http.Response)
	if resp.StatusCode == 200 {
		return nil
	}
	return errors.NewError("internal error", resp.StatusCode)
}
