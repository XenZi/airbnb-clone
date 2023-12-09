package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"notifications-service/domains"
	"notifications-service/services"
	"notifications-service/utils"
)


type NotificationHandler struct {
	service *services.NotificationService
}

func NewNotificationHandler(service *services.NotificationService) *NotificationHandler{
	return &NotificationHandler{
		service: service,
	}
}

func (nh NotificationHandler) CreateNotification(rw http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var notificationData domains.UserNotification
	if err := decoder.Decode(&notificationData); err != nil {
		utils.WriteErrorResp(err.Error(), 500, "api/login", rw)
		return
	}
	log.Println(notificationData)
}