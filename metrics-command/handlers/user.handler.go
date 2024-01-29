package handlers

import (
	"encoding/json"
	"log"
	"metrics-command/commands/handler"
	"metrics-command/commands/user_joined"
	"metrics-command/commands/user_left"
	"metrics-command/domains"
	"metrics-command/utils"
	"net/http"
)

type UserHandler struct {
	handler handler.Handler
}

func NewUserHandler(handler handler.Handler) *UserHandler {
	return &UserHandler{
		handler: handler,
	}
}

func (h UserHandler) CreateJoinedAt(w http.ResponseWriter, r *http.Request) {
	var req domains.UserJoined
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		utils.WriteErrorResp(err.Error(), 400, "api/metrics/joinedAt", w)
		return
	}

	command := user_joined.NewCommand(req.UserID, req.AccommodationID, req.JoinedAt)
	err = h.handler.Handle(command)
	if err != nil {
		log.Println(err)
		utils.WriteErrorResp(err.Error(), 400, "api/metrics/joinedAt", w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h UserHandler) CreateLeftAt(w http.ResponseWriter, r *http.Request) {
	var req domains.UserLeft
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		utils.WriteErrorResp(err.Error(), 400, "api/metrics/joinedAt", w)
		return
	}

	command := user_left.NewCommand(req.UserID, req.AccommodationID, req.LeftAt)
	err = h.handler.Handle(command)
	if err != nil {
		log.Println(err)
		utils.WriteErrorResp(err.Error(), 400, "api/metrics/joinedAt", w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
