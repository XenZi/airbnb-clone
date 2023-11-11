package handler

import (
	"auth-service/domains"
	"auth-service/errors"
	"auth-service/services"
	"auth-service/utils"
	"context"
	"encoding/json"
	"github.com/XenZi/airbnb-clone/api-gateway/proto/auth_service"
	"net/http"
	"time"
)

type AuthHandler struct {
	UserService *services.UserService
	auth_service.UnimplementedAuthServiceServer
}

func (a AuthHandler) Login(ctx context.Context, request *auth_service.LoginRequestDTO) (*auth_service.MixLoginAndError, error) {
	response := &auth_service.MixLoginAndError{}
	loginUser := domains.LoginUser{Email: request.Email, Password: request.Password}
	successfullyLoggedData, err := a.UserService.LoginUser(loginUser)
	if err != nil {
		errorResponse := &auth_service.BaseErrorHttpResponse{
			Message: err.GetErrorMessage(),
			Status:  int32(err.GetErrorStatus()),
			Path:    "api/login",
			Time:    time.Now().String(),
		}

		response.MixLogin = &auth_service.MixLoginAndError_Err{
			Err: errorResponse,
		}
		return response, nil
	}
	userDTO := &auth_service.UserDTO{
		Id:       successfullyLoggedData.User.ID,
		Username: successfullyLoggedData.User.Username,
		Email:    successfullyLoggedData.User.Email,
		Role:     successfullyLoggedData.User.Role,
	}
	successfullyLoggedUser := &auth_service.SuccessfullyLoggedUser{
		Token: successfullyLoggedData.Token,
		User:  userDTO,
	}
	response = &auth_service.MixLoginAndError{
		MixLogin: &auth_service.MixLoginAndError_Res{
			Res: successfullyLoggedUser,
		},
	}
	return response, nil
}
func (a AuthHandler) RegisterHandler(r http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var registerData domains.RegisterUser
	if err := decoder.Decode(&registerData); err != nil {
		utils.WriteErrorResp(errors.ErrInternalServerError().Error(), 500, "api/login", r)
	}
	userData, err := a.UserService.CreateUser(registerData)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/register", r)
		return
	}
	utils.WriteResp(userData, 201, r)
}
