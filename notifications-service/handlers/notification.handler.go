package handlers

import (
	"net/http"
	"notifications-service/services"
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

}