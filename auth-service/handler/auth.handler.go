package handler

import (
	"auth-service/config"
	"auth-service/domains"
	"auth-service/services"
	"auth-service/utils"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/trace"
)

type AuthHandler struct {
	UserService *services.UserService
	Tracer      trace.Tracer
	Logger      *config.Logger
}

func (a AuthHandler) LoginHandler(rw http.ResponseWriter, h *http.Request) {
	ctx, span := a.Tracer.Start(h.Context(), "AuthHandler.Login")
	defer span.End()
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var loginData domains.LoginUser
	if err := decoder.Decode(&loginData); err != nil {
		utils.WriteErrorResp(err.Error(), 500, "api/login", rw)
		return
	}
	jwtToken, err := a.UserService.LoginUser(ctx, loginData)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/login", rw)
		return
	}
	utils.WriteResp(jwtToken, 200, rw)
}

func (a AuthHandler) RegisterHandler(r http.ResponseWriter, h *http.Request) {
	ctx, cancel := context.WithTimeout(h.Context(), time.Second*3)
	defer cancel()
	ctx, span := a.Tracer.Start(ctx, "AuthHandler.Register")
	defer span.End()
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var registerData domains.RegisterUser
	if err := decoder.Decode(&registerData); err != nil {
		utils.WriteErrorResp("Internal server error", 500, "api/login", r)
	}
	userData, err := a.UserService.CreateUser(ctx, registerData)
	log.Println("E$RROR IN HANDL:ER", err)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/register", r)
		return
	}
	utils.WriteResp(userData, 201, r)
}

func (a AuthHandler) ConfirmAccount(r http.ResponseWriter, h *http.Request) {
	ctx, cancel := context.WithTimeout(h.Context(), time.Second*3)
	defer cancel()
	ctx, span := a.Tracer.Start(ctx, "AuthHandler.ConfirmAccount")
	defer span.End()
	vars := mux.Vars(h)
	token := vars["token"]
	if token == "" {
		utils.WriteErrorResp("Bad request", 400, "api/confirm-account", r)
		return
	}
	user, err := a.UserService.ConfirmUserAccount(ctx, token)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/confirm-account", r)
		return
	}
	utils.WriteResp(user, 200, r)
}

func (a AuthHandler) RequestResetPassword(r http.ResponseWriter, h *http.Request) {
	ctx, cancel := context.WithTimeout(h.Context(), time.Second*3)
	defer cancel()
	ctx, span := a.Tracer.Start(ctx, "AuthHandler.RequestResetPassword")
	defer span.End()
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var requestData domains.RequestResetPassword
	if err := decoder.Decode(&requestData); err != nil {
		utils.WriteErrorResp(err.Error(), 500, "api/request-reset-password", r)
		return
	}
	res, err := a.UserService.RequestResetPassword(ctx, requestData.Email)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/request-reset-password", r)
		return
	}
	utils.WriteResp(res, 200, r)
}

func (a AuthHandler) ResetPassword(r http.ResponseWriter, h *http.Request) {
	ctx, cancel := context.WithTimeout(h.Context(), time.Second*3)
	defer cancel()
	ctx, span := a.Tracer.Start(ctx, "AuthHandler.ResetPassword")
	defer span.End()
	vars := mux.Vars(h)
	token := vars["token"]
	if token == "" {
		utils.WriteErrorResp("Bad request", 400, "api/confirm-account", r)
		return
	}
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var requestData domains.ResetPassword
	if err := decoder.Decode(&requestData); err != nil {
		utils.WriteErrorResp(err.Error(), 500, "api/request-reset-password", r)
		return
	}
	user, err := a.UserService.ResetPassword(ctx, requestData, token)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/reset-password", r)
		return
	}
	utils.WriteResp(user, 200, r)
}

func (a AuthHandler) ChangePassword(r http.ResponseWriter, h *http.Request) {
	ctx, cancel := context.WithTimeout(h.Context(), time.Second*3)
	defer cancel()
	log.Println("CTX VALUEEEE: ", h.Context().Value("userID"))
	ctx, span := a.Tracer.Start(ctx, "AuthHandler.ChangePassword")
	defer span.End()
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var requestData domains.ChangePassword
	if err := decoder.Decode(&requestData); err != nil {
		utils.WriteErrorResp(err.Error(), 500, "api/change-password", r)
		return
	}
	vars := mux.Vars(h)
	id := vars["id"]
	if id == "" {
		utils.WriteErrorResp("Bad request", 400, "api/confirm-account", r)
		return
	}
	resp, err := a.UserService.ChangePassword(ctx, requestData, id)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/auth/change-password", r)
		return
	}
	utils.WriteResp(resp, 200, r)
}

func (a AuthHandler) UpdateCredentials(r http.ResponseWriter, h *http.Request) {
	ctx, cancel := context.WithTimeout(h.Context(), time.Second*3)
	defer cancel()
	userID := h.Context().Value("userID").(string)
	ctx, span := a.Tracer.Start(ctx, "AuthHandler.UpdateCredentials")
	defer span.End()
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var requestData domains.User
	if err := decoder.Decode(&requestData); err != nil {
		utils.WriteErrorResp(err.Error(), 500, "api/change-password", r)
		return
	}
	res, err := a.UserService.UpdateCredentials(ctx, userID, requestData)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/update-credentials", r)
		return
	}
	utils.WriteResp(res, 201, r)
}

func (a AuthHandler) DeleteUser(r http.ResponseWriter, h *http.Request) {
	ctx, cancel := context.WithTimeout(h.Context(), time.Second*3)
	defer cancel()
	ctx, span := a.Tracer.Start(ctx, "AuthHandler.DeleteUser")
	defer span.End()
	vars := mux.Vars(h)
	id := vars["id"]
	if id == "" {
		utils.WriteErrorResp("Bad request", 400, "api/delete-user", r)
		return
	}
	resp, err := a.UserService.DeleteUserById(ctx, id)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/delete", r)
		return
	}
	utils.WriteResp(resp, 200, r)
}

func (a AuthHandler) All(r http.ResponseWriter, h *http.Request) {
	utils.WriteResp("Cool", 200, r)
}
