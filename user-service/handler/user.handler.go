package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"user-service/domain"
	"user-service/service"
	"user-service/utils"
)

type UserHandler struct {
	UserService *service.UserService
}

func (u UserHandler) CreateHandler(rw http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var createData domain.CreateUser
	if err := decoder.Decode(&createData); err != nil {
		utils.WriteErrorResponse(err.Error(), 500, "api/create", rw)
		return
	}
	user, err := u.UserService.CreateUser(createData)
	if err != nil {
		utils.WriteErrorResponse(err.GetErrorMessage(), err.GetErrorStatus(), "api/create", rw)
		return
	}
	utils.WriteResp(user, 200, rw)
}

func (u UserHandler) UpdateHandler(rw http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var updateData domain.CreateUser
	if err := decoder.Decode(&updateData); err != nil {
		utils.WriteErrorResponse(err.Error(), 500, "api/update", rw)
		return
	}
	user, err := u.UserService.UpdateUser(updateData)
	if err != nil {
		utils.WriteErrorResponse(err.GetErrorMessage(), err.GetErrorStatus(), "api/update", rw)
		return
	}
	utils.WriteResp(user, 200, rw)
}

func (u UserHandler) CredsHandler(rw http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var updateData domain.CreateUser
	if err := decoder.Decode(&updateData); err != nil {
		utils.WriteErrorResponse(err.Error(), 500, "api/update", rw)
		return
	}
	user, err := u.UserService.UpdateUserCreds(updateData)
	if err != nil {
		utils.WriteErrorResponse(err.GetErrorMessage(), err.GetErrorStatus(), "api/update", rw)
		return
	}
	utils.WriteResp(user, 200, rw)
}

func (u UserHandler) GetAllHandler(rw http.ResponseWriter, h *http.Request) {
	userCollection, err := u.UserService.GetAllUsers()
	if err != nil {
		utils.WriteErrorResponse(err.GetErrorMessage(), err.GetErrorStatus(), "api/get-all", rw)
		return
	}
	utils.WriteResp(userCollection, 200, rw)
}

func (u UserHandler) GetUserById(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]
	user, err := u.UserService.GetUserById(id)
	if err != nil {
		utils.WriteErrorResponse(err.GetErrorMessage(), err.GetErrorStatus(), "api/get-user", rw)
		return
	}
	utils.WriteResp(user, 200, rw)
}

func (u UserHandler) DeleteHandler(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]
	err := u.UserService.DeleteUser(id)
	if err != nil {
		utils.WriteErrorResponse(err.GetErrorMessage(), err.GetErrorStatus(), "api/delete", rw)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusNoContent)

}
