package client

import (
	"accommodations-service/config"
	"accommodations-service/domain"
	"accommodations-service/errors"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sony/gobreaker"
	"log"
	"net/http"
)

type UserClient struct {
	address        string
	client         *http.Client
	circuitBreaker *gobreaker.CircuitBreaker
	logger         *config.Logger
}
type HostUser struct {
	ID            string  `bson:"_id,omitempty" json:"id"`
	Username      string  `json:"username" bson:"username"`
	Email         string  `json:"email" bson:"email"`
	Role          string  `json:"role" bson:"role"`
	FirstName     string  `json:"firstName" bson:"firstName"`
	LastName      string  `json:"lastName" bson:"lastName"`
	Residence     string  `json:"residence" bson:"residence"`
	Age           int     `json:"age" bson:"age"`
	Rating        float64 `json:"rating"`
	Distinguished bool    `json:"distinguished"`
}

func NewUserClient(host, port string, client *http.Client, circuitBreaker *gobreaker.CircuitBreaker, logger *config.Logger) *UserClient {
	return &UserClient{
		address:        fmt.Sprintf("http://%s:%s", host, port),
		client:         http.DefaultClient,
		circuitBreaker: circuitBreaker,
		logger:         logger,
	}
}

func (uc UserClient) GetUserById(ctx context.Context, id string) (*HostUser, *errors.ErrorStruct) {
	cbResp, err := uc.circuitBreaker.Execute(func() (interface{}, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, uc.address+"/"+id, nil)
		if err != nil {
			uc.logger.LogError("accommodations-client", fmt.Sprintf("Unable to send request to user service"))
			uc.logger.LogError("accommodation-client", fmt.Sprintf("Error:"+err.Error()))
			log.Println(err)
			return nil, err
		}
		uc.logger.LogInfo("accommodation-client", fmt.Sprintf("Successfully sent request to user service"))
		return uc.client.Do(req)
	})

	if err != nil {
		uc.logger.LogError("accommodations-client", fmt.Sprintf("Unable to send request to user service"))
		uc.logger.LogError("accommodation-client", fmt.Sprintf("Error:"+err.Error()))
		log.Println("ERR FROM GGG ", err)
		return nil, errors.NewError("Nothing to parse", 500)
	}

	resp := cbResp.(*http.Response)
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		baseResp := domain.BaseHttpResponse{}
		err := json.NewDecoder(resp.Body).Decode(&baseResp)
		if err != nil {
			uc.logger.LogError("accommodations-client", fmt.Sprintf("Unable to decode base response"))
			uc.logger.LogError("accommodation-client", fmt.Sprintf("Error:"+err.Error()))
			return nil, errors.NewError(err.Error(), 500)
		}
		log.Println("BASE RESP JE", baseResp)

		// Check if the data is a map
		if userData, ok := baseResp.Data.(map[string]interface{}); ok {
			// Create a HostUser manually
			user := HostUser{
				ID:            userData["id"].(string),
				FirstName:     userData["firstName"].(string),
				LastName:      userData["lastName"].(string),
				Age:           int(userData["age"].(float64)),
				Distinguished: userData["distinguished"].(bool),
				Email:         userData["email"].(string),
				Rating:        float64(int(userData["rating"].(float64))),
				Residence:     userData["residence"].(string),
				Role:          userData["role"].(string),
				Username:      userData["username"].(string),
			}
			uc.logger.LogInfo("accommodation-client", fmt.Sprintf("Successfully sent request to user service and retrieved data"))
			log.Println("User je", user)
			return &user, nil
		} else {
			uc.logger.LogError("accommodations-client", fmt.Sprintf("Unable to map data to the object"))
			uc.logger.LogError("accommodation-client", fmt.Sprintf("Error:"+err.Error()))
			return nil, errors.NewError("Invalid data type in response", 500)
		}
	}

	baseResp := domain.BaseErrorHttpResponse{}
	err = json.NewDecoder(resp.Body).Decode(&baseResp)
	if err != nil {
		uc.logger.LogError("accommodations-client", fmt.Sprintf("Unable to decode response body"))
		uc.logger.LogError("accommodation-client", fmt.Sprintf("Error:"+err.Error()))
		return nil, errors.NewError(err.Error(), 500)
	}

	log.Println(baseResp)
	log.Println(baseResp.Error)
	return nil, errors.NewError(baseResp.Error, baseResp.Status)
}
