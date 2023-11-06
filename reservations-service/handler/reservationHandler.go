package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"reservation-service/domain"
	"reservation-service/repository"
)

type KeyProduct struct{}

type ReservationHandler struct {
	logger *log.Logger

	repo *repository.ReservationRepo
}

func NewReservationsHandler(l *log.Logger, r *repository.ReservationRepo) *ReservationHandler {
	return &ReservationHandler{l, r}
}

func (r *ReservationHandler) CreateReservationById(rw http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)

	var loginData domain.Reservation
	if err := decoder.Decode(&loginData); err != nil {

		return
	}

	reservationById, err := r.repo.InsertReservationById(&loginData)
	r.logger.Println(reservationById)
	if err != nil {
		r.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusCreated)
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
