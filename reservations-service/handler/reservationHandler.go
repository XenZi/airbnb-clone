package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"reservation-service/domain"
	"reservation-service/service"
	"reservation-service/utils"

	"github.com/gorilla/mux"
)

type KeyProduct struct{}

type ReservationHandler struct {
	logger *log.Logger

	ReservationService *service.ReservationService
}

func NewReservationsHandler(l *log.Logger, rs *service.ReservationService) *ReservationHandler {
	return &ReservationHandler{l, rs}
}

func (r *ReservationHandler) CreateReservation(rw http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var res domain.Reservation
	if err := decoder.Decode(&res); err != nil {
		utils.WriteErrorResp(err.Error(), 500, "api/reservations", rw)
		return
	}
	newRes, err := r.ReservationService.CreateReservation(res)
	if err != nil {
		utils.WriteErrorResp(err.Message, err.Status, "api/reservations", rw)
		return
	}
	utils.WriteResp(newRes, 201, rw)
}
func (r *ReservationHandler) CreateAvailability(rw http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var avl domain.FreeReservation
	if err := decoder.Decode(&avl); err != nil {
		utils.WriteErrorResp(err.Error(), 500, "api/availability", rw)
		return
	}
	log.Println("USLO U CREATE")
	newAvl, err := r.ReservationService.CreateAvailability(avl)
	if err != nil {
		utils.WriteErrorResp(err.Message, err.Status, "api/availability", rw)
		return
	}
	utils.WriteResp(newAvl, 201, rw)

}
func (rh *ReservationHandler) GetReservationsByUser(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	reservations, err := rh.ReservationService.GetReservationsByUser(userID)
	if err != nil {
		utils.WriteErrorResp(err.Message, err.Status, "api/reservations/user/{userId}", rw)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(reservations)
}
func (rh ReservationHandler) GetReservationsByAccommodation(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accommodationID := vars["accommodationID"]

	reservations, err := rh.ReservationService.GetReservationsByAccommodation(accommodationID)
	if err != nil {
		utils.WriteErrorResp(err.Message, err.Status, "api/accommodation/sreservations/{accommodationID}", rw)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(reservations)
}

func (rh *ReservationHandler) DeleteReservationById(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	country := vars["country"]

	deletedReservation, err := rh.ReservationService.DeleteReservationById(country, id)
	if err != nil {
		utils.WriteErrorResp(err.Message, err.Status, "api/reservations/{country}/{id}", rw)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(deletedReservation)
}
