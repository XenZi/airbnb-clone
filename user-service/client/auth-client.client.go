package client

import (
	"context"
	"fmt"
	"github.com/sony/gobreaker"
	"log"
	"net/http"
	"user-service/errors"
)

type AuthClient struct {
	address        string
	client         *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
}

func NewAuthClient(host, port string, client *http.Client, circuitBreaker *gobreaker.CircuitBreaker) *AuthClient {
	return &AuthClient{
		address:        fmt.Sprintf("http://%s:%s", host, port),
		client:         client,
		circuitBreaker: circuitBreaker,
	}
}

func (ac AuthClient) DeleteUserAuth(ctx context.Context, id string) *errors.ErrorStruct {
	log.Println("Poslato u bezdan")
	cbResp, err := ac.circuitBreaker.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, ac.address+"/"+id, http.NoBody)
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
