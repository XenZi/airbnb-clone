package handlers

import (
	"encoding/json"
	"log"
	"metrics-command/commands/handler"
	"metrics-command/commands/user_reserved"
	"metrics-command/domains"
	"metrics-command/utils"
	"net/http"
)

type ReservationHandler struct {
	handler handler.Handler
}

func NewReservationHandler(handler handler.Handler) *ReservationHandler {
	return &ReservationHandler{
		handler: handler,
	}
}

func (h ReservationHandler) CreateReserved(w http.ResponseWriter, r *http.Request) {
	var req domains.Reservation
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		utils.WriteErrorResp(err.Error(), 400, "api/metrics/joinedAt", w)
		return
	}
	command := user_reserved.NewCommand(req.UserID, req.AccommodationID, req.ReservedAt)
	err = h.handler.Handle(command)
	if err != nil {
		log.Println(err)
		utils.WriteErrorResp(err.Error(), 400, "api/metrics/joinedAt", w)
		return
	}

	utils.WriteResp(string("Successfully inserted"), 200, w)

}
