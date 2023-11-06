package handlers

import (
	"accommodations-service/domain"
	"accommodations-service/repository"
	"context"
	"log"
	"net/http"
)

type KeyProduct struct{}

type AccommodationsHandler struct {
	logger *log.Logger

	repo *repository.AccommodationRepo
}

func NewAccommodationsHandler(l *log.Logger, r *repository.AccommodationRepo) *AccommodationsHandler {
	return &AccommodationsHandler{l, r}
}

func (a *AccommodationsHandler) CreateAccommodationById(rw http.ResponseWriter, h *http.Request) {
	accommodationById := h.Context().Value(KeyProduct{}).(*domain.Accommodation)
	err := a.repo.InsertAccommodationById(accommodationById)
	if err != nil {
		a.logger.Print("Database exception: ", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

func (a *AccommodationsHandler) MiddlewareAccommodationByIdDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		patient := &domain.Accommodation{}
		err := patient.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			a.logger.Fatal(err)
			return
		}
		ctx := context.WithValue(h.Context(), KeyProduct{}, patient)
		h = h.WithContext(ctx)
		next.ServeHTTP(rw, h)
	})
}
