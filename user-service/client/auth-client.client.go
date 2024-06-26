package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sony/gobreaker"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"user-service/domain"
	"user-service/errors"
)

type AuthClient struct {
	address        string
	client         *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
	tracer         trace.Tracer
}

func NewAuthClient(host, port string, client *http.Client, circuitBreaker *gobreaker.CircuitBreaker, tracer trace.Tracer) *AuthClient {
	return &AuthClient{
		address:        fmt.Sprintf("http://%s:%s", host, port),
		client:         client,
		circuitBreaker: circuitBreaker,
		tracer:         tracer,
	}
}

func (ac AuthClient) DeleteUserAuth(ctx context.Context, id string) *errors.ErrorStruct {
	ctx, span := ac.tracer.Start(ctx, "AuthnClient.DeleteUserAuth")
	defer span.End()
	cbResp, err := ac.circuitBreaker.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, ac.address+"/"+id, http.NoBody)
		if err != nil {
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
