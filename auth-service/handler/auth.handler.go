package handler

import (
	"auth-service/domains"
	"auth-service/services"
	"auth-service/utils"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type AuthHandler struct {
	UserService *services.UserService
}

func (a AuthHandler) LoginHandler(rw http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var loginData domains.LoginUser
	if err := decoder.Decode(&loginData); err != nil {
		utils.WriteErrorResp(err.Error(), 500, "api/login", rw)
		return
	}
	jwtToken, err := a.UserService.LoginUser(loginData)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/login", rw)
		return
	}
	utils.WriteResp(jwtToken, 200, rw)
}

func (a AuthHandler) RegisterHandler(r http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var registerData domains.RegisterUser
	if err := decoder.Decode(&registerData); err != nil {
		utils.WriteErrorResp("Internal server error", 500, "api/login", r)
	}
	userData, err := a.UserService.CreateUser(registerData)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/register", r)
		return
	}
	utils.WriteResp(userData, 201, r)
}

func (a AuthHandler) ConfirmAccount(r http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	token := vars["token"]
	if token == "" {
		utils.WriteErrorResp("Bad request", 400, "api/confirm-account", r)
		return
	}
	user, err := a.UserService.ConfirmUserAccount(token)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/confirm-account",r)
		return
	}
	utils.WriteResp(user, 200, r)
}

func (a AuthHandler) RequestResetPassword(r http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var requestData domains.RequestResetPassword
	if err := decoder.Decode(&requestData); err != nil {
		utils.WriteErrorResp(err.Error(), 500, "api/request-reset-password", r)
		return
	}
	res, err := a.UserService.RequestResetPassword(requestData.Email)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/request-reset-password", r)
		return
	}
	utils.WriteResp(res, 200, r)
}

func (a AuthHandler) ResetPassword (r http.ResponseWriter, h  *http.Request) {
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
	user, err := a.UserService.ResetPassword(requestData, token)
	if err != nil { 
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/reset-password", r)
		return
	}
	utils.WriteResp(user, 200, r)
}

func (a AuthHandler) ChangePassword(r http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var requestData domains.ChangePassword
	if err := decoder.Decode(&requestData); err != nil {
		utils.WriteErrorResp(err.Error(), 500, "api/change-password", r)
		return
	}
	userID := h.Context().Value("userID").(string)
	resp, err := a.UserService.ChangePassword(requestData, userID)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/auth/change-password", r)
		return
	}
	utils.WriteResp(resp, 200, r)
}