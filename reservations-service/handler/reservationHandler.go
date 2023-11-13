package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"reservation-service/domain"
	"reservation-service/repository"

	"github.com/gorilla/mux"
)

type KeyProduct struct{}

type ReservationHandler struct {
	logger *log.Logger

	repo *repository.ReservationRepo
}

func NewReservationsHandler(l *log.Logger, r *repository.ReservationRepo) *ReservationHandler {
	return &ReservationHandler{l, r}
}

func (r *ReservationHandler) CreateReservationByUser(rw http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	var reservation domain.Reservation
	if err := decoder.Decode(&reservation); err != nil {
		log.Println(err)
		return
	}
	reservationById, err := r.repo.InsertReservationByUser(&reservation)
	r.logger.Println(reservationById)
	if err != nil {
		r.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}
func (r *ReservationHandler) CreateReservationByAccommodation(rw http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	var reservation domain.Reservation
	if err := decoder.Decode(&reservation); err != nil {
		log.Println(err)
		return
	}
	reservationById, err := r.repo.InsertReservationByAccommodantion(&reservation)
	r.logger.Println(reservationById)
	if err != nil {
		r.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

func (r *ReservationHandler) GetReservationsByUser(rw http.ResponseWriter, req *http.Request) {
	userID := mux.Vars(req)["userId"]

	reservations, err := r.repo.GetReservationsByUser(userID)
	if err != nil {
		http.Error(rw, "Failed to get reservations", http.StatusInternalServerError)
		return
	}

	if reservations == nil {
		http.Error(rw, "No reservations found for the user", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(rw).Encode(reservations)
	if err != nil {
		http.Error(rw, "Unable to convert to JSON", http.StatusInternalServerError)
		r.logger.Fatal("Unable to convert to JSON: ", err)
		return
	}
}

/*func (r *ReservationHandler) MiddlewareReservationByIdDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		patient := &domain.Reservation{}
		err := patient.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			r.logger.Fatal(err)
			return
		}
		ctx := context.WithValue(h.Context(), KeyProduct{}, patient)
		h = h.WithContext(ctx)
		next.ServeHTTP(rw, h)
	})
}
*/
