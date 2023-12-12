package handlers

import (
	"encoding/json"
	"net/http"
	"notifications-service/domains"
	"notifications-service/services"
	"notifications-service/utils"

	"github.com/gorilla/mux"
)


type NotificationHandler struct {
	service *services.NotificationService
}

func NewNotificationHandler(service *services.NotificationService) *NotificationHandler{
	return &NotificationHandler{
		service: service,
	}
}

func (nh NotificationHandler) CreateNewUserNotification(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]
	if id == "" {
		utils.WriteErrorResp("Bad request", 400, "api/notifications", rw)
		return
	}
	resp, err := nh.service.CreateNewUserNotification(id)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/notifications", rw)
		return
	}
	utils.WriteResp(resp, 201, rw)
}

func (nh NotificationHandler) CreateNewNotificationForUser(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]
	if id == "" {
		utils.WriteErrorResp("Bad request", 400, "api/notifications", rw)
		return
	}
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var requestData domains.Notification
	if err := decoder.Decode(&requestData); err != nil {
		utils.WriteErrorResp(err.Error(), 500, "api/notifications", rw)
		return
	}
	resp, err := nh.service.PushNewNotificationToUser(id, requestData)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/notifications", rw)
		return
	}
	utils.WriteResp(resp, 201, rw)
}

func (nh NotificationHandler) ReadAllNotifications(rw http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var requestData domains.UserNotificationDTO
	if err := decoder.Decode(&requestData); err != nil {
		utils.WriteErrorResp(err.Error(), 500, "api/notifications", rw)
		return
	}
	resp, err := nh.service.ReadAllNotifications(requestData)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/notifications",rw)
		return
	}
	utils.WriteResp(resp, 200, rw)
}

func (nh NotificationHandler) GetAllNotificationsByID(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]
	if id == "" {
		utils.WriteErrorResp("Bad request", 400, "api/notifications", rw)
		return
	}
	resp, err := nh.service.FindAllNotificationsByID(id)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/notifications",rw)
		return
	}
	utils.WriteResp(resp, 200, rw) 
}