package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"reservation-service/domain"
	"reservation-service/service"
	"reservation-service/utils"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/trace"
)

type KeyProduct struct{}

type ReservationHandler struct {
	logger             *log.Logger
	ReservationService *service.ReservationService
	Tracer             trace.Tracer
}

func NewReservationsHandler(l *log.Logger, rs *service.ReservationService, tr trace.Tracer) *ReservationHandler {
	return &ReservationHandler{l, rs, tr}
}

func (r *ReservationHandler) CreateReservation(rw http.ResponseWriter, h *http.Request) {
	ctx, span := r.Tracer.Start(h.Context(), "ReservationHandler.CreateReservation")
	defer span.End()
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
	accommodationID := vars["accommodationId"]

	avl, err := rh.ReservationService.GetAvailabilityForAccommodation(accommodationID)
	if err != nil {
		utils.WriteErrorResp(err.Message, err.Status, "api/{accommodationId}/availability", rw)
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

	reservations, erro := rh.ReservationService.ProcessDateRange(request.AccommodationIDs, request.DateRange)
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
	endDate := vars["endDate"]

	deletedReservation, err := rh.ReservationService.DeleteReservationById(country, id, userID, hostID, accommodationID, endDate)
	if err != nil {
		utils.WriteErrorResp(err.Message, err.Status, "api/reservations/{country}/{id}/{userID}/{hostID}/{accommodationID}/{endDate}", rw)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(deletedReservation)
}

func (rh *ReservationHandler) GetCancelationPercentage(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hostID := vars["hostID"]
	percentage, err := rh.ReservationService.CalculatePercentageCanceled(hostID)
	if err != nil {
		utils.WriteErrorResp(err.Message, err.Status, "/api/reservations/percentage-cancelation/{hostID}", rw)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(percentage)
}
func (rh *ReservationHandler) GetReservationsByAccommodationWithEndDate(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accommodationID := vars["accommodationId"]
	userID := vars["userId"]

	reservations, err := rh.ReservationService.GetReservationsByAccommodationWithEndDate(accommodationID, userID)
	if err != nil {
		utils.WriteErrorResp(err.Message, err.Status, "api/reservations/{accommodationId}/{userId}", rw)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	utils.WriteResp(reservations, 200, rw)
}
func (rh *ReservationHandler) GetReservationsByHostWithEndDate(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hostID := vars["hostId"]
	userID := vars["userId"]

	reservations, err := rh.ReservationService.GetReservationsByHostWithEndDate(hostID, userID)
	if err != nil {
		utils.WriteErrorResp(err.Message, err.Status, "api/reservations/{hostId}/{userId}", rw)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	utils.WriteResp(reservations, 200, rw)
}

func (rh *ReservationHandler) UpdateAvailability(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accommodationID := vars["accommodationId"]
	id := vars["id"]
	country := vars["country"]
	priceStr := vars["price"]

	price, err := strconv.Atoi(priceStr)
	if err != nil {
		http.Error(w, "Invalid price parameter", http.StatusBadRequest)
		return
	}

	var updatedReservation domain.FreeReservation
	err = json.NewDecoder(r.Body).Decode(&updatedReservation)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, reservationErr := rh.ReservationService.UpdateAvailability(accommodationID, id, country, price, &updatedReservation)
	if reservationErr != nil {
		http.Error(w, reservationErr.Message, 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	utils.WriteResp(result, 200, w)
}

func (rh *ReservationHandler) GetAccommodationIDsByMaxPrice(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	maxPriceStr := vars["maxPrice"]
	maxPrice, err := strconv.Atoi(maxPriceStr)
	if err != nil {
		utils.WriteErrorResp(err.Error(), 500, "api/reservations/price/myprice/janko/mateja/aca/{maxPrice}", rw)
		return
	}
	log.Println(maxPrice)
	accommodations, erro := rh.ReservationService.GetAccommodationIDsByMaxPrice(maxPrice)
	if erro != nil {
		utils.WriteErrorResp(err.Error(), 500, "api/reservations/price/myprice/janko/mateja/aca/{maxPrice}", rw)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	utils.WriteResp(accommodations, http.StatusOK, rw)
}
