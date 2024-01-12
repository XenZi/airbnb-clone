package client

import (
	"fmt"
	"net/http"

	"github.com/sony/gobreaker"
)

type UserClient struct {
	address        string
	client         *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
}

func NewUserClient(host, port string, client *http.Client, cb *gobreaker.CircuitBreaker) *UserClient {
	return &UserClient{
		address:        fmt.Sprintf("http://%s:%s"),
		client:         client,
		circuitBreaker: cb,
	}
}
