package handler

import (
	"encoding/json"
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
	}
	utils.WriteResp(user, 200, rw)
}
