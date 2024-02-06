package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"notifications-service/domains"
	"notifications-service/services"
	"notifications-service/utils"
	"time"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/trace"
)

type NotificationHandler struct {
	service *services.NotificationService
	tracer  trace.Tracer
}

func NewNotificationHandler(service *services.NotificationService, tracer trace.Tracer) *NotificationHandler {
	return &NotificationHandler{
		service: service,
		tracer:  tracer,
	}
}

func (nh NotificationHandler) CreateNewUserNotification(rw http.ResponseWriter, h *http.Request) {
	ctx, span := nh.tracer.Start(h.Context(), "NotificationHandler.CreateNewUserNotification")
	defer span.End()

	vars := mux.Vars(h)
	id := vars["id"]
	if id == "" {
		utils.WriteErrorResp("Bad request", 400, "api/notifications", rw)
		return
	}
	resp, err := nh.service.CreateNewUserNotification(ctx, id)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/notifications", rw)
		return
	}
	utils.WriteResp(resp, 201, rw)
}

func (nh NotificationHandler) CreateNewNotificationForUser(rw http.ResponseWriter, h *http.Request) {
	ctx, cancel := context.WithTimeout(h.Context(), 10*time.Second)
	defer cancel()
	ctx, span := nh.tracer.Start(ctx, "NotificationHandler.CreateNewUserNotification")
	defer span.End()
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

	resp, err := nh.service.PushNewNotificationToUser(ctx, id, requestData)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/notifications", rw)
		return
	}
	utils.WriteResp(resp, 201, rw)
}

func (nh NotificationHandler) ReadAllNotifications(rw http.ResponseWriter, h *http.Request) {
	ctx, span := nh.tracer.Start(h.Context(), "NotificationHandler.CreateNewUserNotification")
	defer span.End()
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var requestData domains.UserNotificationDTO
	if err := decoder.Decode(&requestData); err != nil {
		utils.WriteErrorResp(err.Error(), 500, "api/notifications", rw)
		return
	}
	resp, err := nh.service.ReadAllNotifications(ctx, requestData)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/notifications", rw)
		return
	}
	utils.WriteResp(resp, 200, rw)
}

func (nh NotificationHandler) GetAllNotificationsByID(rw http.ResponseWriter, h *http.Request) {
	ctx, span := nh.tracer.Start(h.Context(), "NotificationHandler.CreateNewUserNotification")
	defer span.End()
	vars := mux.Vars(h)
	id := vars["id"]
	if id == "" {
		utils.WriteErrorResp("Bad request", 400, "api/notifications", rw)
		return
	}
	resp, err := nh.service.FindAllNotificationsByID(ctx, id)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/notifications", rw)
		return
	}
	utils.WriteResp(resp, 200, rw)
}
