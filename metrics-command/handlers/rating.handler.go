package handlers

import (
	"encoding/json"
	"log"
	"metrics-command/commands/handler"
	"metrics-command/commands/user_rated"
	"metrics-command/domains"
	"metrics-command/utils"
	"net/http"
)

type RatingHandler struct {
	handler handler.Handler
}

func NewRatingHandler(handler handler.Handler) *RatingHandler {
	return &RatingHandler{
		handler: handler,
	}
}

func (h RatingHandler) CreateRatedAt(w http.ResponseWriter, r *http.Request) {
	var req domains.UserRate
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		utils.WriteErrorResp(err.Error(), 400, "api/metrics/joinedAt", w)
		return
	}
	command := user_rated.NewCommand(req.UserID, req.AccommodationID, req.RatedAt)
	err = h.handler.Handle(command)
	if err != nil {
		log.Println(err)
		utils.WriteErrorResp(err.Error(), 400, "api/metrics/joinedAt", w)
		return
	}
	utils.WriteResp(string("Successfully rated"), 200, w)
}
