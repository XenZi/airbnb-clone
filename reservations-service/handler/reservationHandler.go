package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"reservation-service/domain"
	"reservation-service/service"
	"reservation-service/utils"
	"time"

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
	ctx, cancel := context.WithTimeout(h.Context(), time.Second*5)
	defer cancel()
	newRes, err := r.ReservationService.CreateReservation(res, ctx)
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
		log.Println("PRvi ErrOR")
		utils.WriteErrorResp(err.Error(), 500, "api/availability", rw)
		return
	}
	log.Println("USLO U CREATE")
	newAvl, err := r.ReservationService.CreateAvailability(avl)
	if err != nil {
		log.Println("DRugI erROr")
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
		utils.WriteErrorResp(err.Message, err.Status, "api/reservations/user/guest/{userId}", rw)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	utils.WriteResp(reservations, 200, rw)
}
func (rh *ReservationHandler) GetReservationsByHost(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hostID := vars["hostId"]

	reservations, err := rh.ReservationService.GetReservationsByUser(hostID)
	if err != nil {
		utils.WriteErrorResp(err.Message, err.Status, "api/reservations/user/{hostId}", rw)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(reservations)
}

func (rh *ReservationHandler) GetAvailabilityForAccommodation(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accommodationID := vars["accommodationID"]

	avl, err := rh.ReservationService.GetAvailabilityForAccommodation(accommodationID)
	if err != nil {
		utils.WriteErrorResp(err.Message, err.Status, "api/{accommodationID}/availability", rw)
	}
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(avl)

}

/*
	func (rh ReservationHandler) GetReservationsByAccommodation(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accommodationID := vars["accommodationID"]

		reservations, err := rh.ReservationService.GetReservationsByAccommodation(accommodationID)
		if err != nil {
			utils.WriteErrorResp(err.Message, err.Status, "api/sreservations/accommodation/{accommodationID}", rw)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(reservations)
	}
*/
func (rh ReservationHandler) GetAvailableDates(rw http.ResponseWriter, r *http.Request) {
	var request domain.CheckAvailabilityRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.WriteErrorResp(err.Error(), 500, "api/reservations/accommodation/dates", rw)
		return
	}
	avl, erro := rh.ReservationService.GetAvailableDates(request.AccommodationID, request.DateRange)
	if erro != nil {
		utils.WriteErrorResp(erro.Error(), 500, "api/reservations/accommodation/dates", rw)
		return
	}
	utils.WriteResp(avl, 201, rw)

}
func (rh ReservationHandler) ReservationsInDateRangeHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.ReservationsInDateRangeRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.WriteErrorResp(err.Error(), 500, "api/reservations/accommodations", w)
		return
	}

	reservations, erro := rh.ReservationService.ReservationsInDateRange(request.AccommodationIDs, request.DateRange)
	log.Println(reservations)
	if erro != nil {
		utils.WriteErrorResp(erro.Error(), 500, "api/reservations/accommodations", w)
		return
	}
	utils.WriteResp(reservations, 201, w)
}

func (rh *ReservationHandler) DeleteReservationById(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	country := vars["country"]
	userID := vars["userID"]
	hostID := vars["hostID"]
	accommodationID := vars["accommodationID"]

	deletedReservation, err := rh.ReservationService.DeleteReservationById(country, id, userID, hostID, accommodationID)
	if err != nil {
		utils.WriteErrorResp(err.Message, err.Status, "api/reservations/{country}/{id}/{userID}/{hostID}/{accommodationID}", rw)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(deletedReservation)
}
