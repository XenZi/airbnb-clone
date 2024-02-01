package client

import (
	"accommodations-service/domain"
	"accommodations-service/errors"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sony/gobreaker"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

type UserClient struct {
	address        string
	client         *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
}
type HostUser struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username      string             `json:"username" bson:"username"`
	Email         string             `json:"email" bson:"email"`
	Role          string             `json:"role" bson:"role"`
	FirstName     string             `json:"firstName" bson:"firstName"`
	LastName      string             `json:"lastName" bson:"lastName"`
	Residence     string             `json:"residence" bson:"residence"`
	Age           int                `json:"age" bson:"age"`
	Rating        float64            `json:"rating"`
	Distinguished bool               `json:"distinguished"`
}

func NewUserClient(host, port string, client *http.Client, circuitBreaker *gobreaker.CircuitBreaker) *UserClient {
	return &UserClient{
		address:        fmt.Sprintf("http://%s:%s", host, port),
		client:         http.DefaultClient,
		circuitBreaker: circuitBreaker,
	}
}

func (uc UserClient) GetUserById(ctx context.Context, id string) (*HostUser, *errors.ErrorStruct) {
	cbResp, err := uc.circuitBreaker.Execute(func() (interface{}, error) {

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, uc.address+"/"+id, nil)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return uc.client.Do(req)
	})

	if err != nil {
		log.Println("ERR FROM GGG ", err)
		return nil, errors.NewError("Nothing to parse", 500)
	}

	resp := cbResp.(*http.Response)
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {

		baseResp := domain.BaseHttpResponse{}
		err := json.NewDecoder(resp.Body).Decode(&baseResp)
		if err != nil {
			return nil, errors.NewError(err.Error(), 500)
		}
		log.Println("BASE RESP JE", baseResp)
		userFromResp := baseResp.Data.(HostUser)

		log.Println("User je", userFromResp)
		return &userFromResp, nil
	}

	baseResp := domain.BaseErrorHttpResponse{}
	err = json.NewDecoder(resp.Body).Decode(&baseResp)
	if err != nil {
		return nil, errors.NewError(err.Error(), 500)
	}

	log.Println(baseResp)
	log.Println(baseResp.Error)
	return nil, errors.NewError(baseResp.Error, baseResp.Status)
}
