package handlers

import (
	"accommodations-service/domain"
	"accommodations-service/errors"
	"accommodations-service/repository"
	"accommodations-service/utils"
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type KeyProduct struct{}

type AccommodationsHandler struct {
	logger *log.Logger

	repo *repository.AccommodationRepo

	validator *utils.Validator
}

func NewAccommodationsHandler(l *log.Logger, r *repository.AccommodationRepo, validator *utils.Validator) *AccommodationsHandler {
	return &AccommodationsHandler{l, r, validator}
}

func (a *AccommodationsHandler) CreateAccommodationById(rw http.ResponseWriter, h *http.Request) {
	accommodationById := h.Context().Value(KeyProduct{}).(*domain.Accommodation)
	a.validator.ValidateAccommodation(accommodationById)
	validatorErrors := a.validator.GetErrors()
	if len(validatorErrors) > 0 {
		var constructedError string
		for _, message := range validatorErrors {
			constructedError += message + "\n"
		}
		http.Error(rw, constructedError, http.StatusBadRequest)
		return
	}
	accommodationById, err := a.repo.InsertAccommodationById(accommodationById)
	if err != nil {
		a.logger.Print("Database exception: ", err)
		utils.WriteErrorResp(err.GetErrorMessage(), 500, "api/accommodations", rw)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (a *AccommodationsHandler) GetAllAccommodations(rw http.ResponseWriter, h *http.Request) {

	accommodations, err := a.repo.GetAllAccommodations()
	if err != nil {
		a.logger.Print("Database exception: ", err)
		utils.WriteErrorResp(err.GetErrorMessage(), 500, "api/accommodations", rw)
	}

	if accommodations == nil {
		return
	}

	if err := accommodations.ToJSON(rw); err != nil {
		a.handleJSONConversionError(errors.NewError("Conversion error", 500), rw)
		return
	}

}

func (a *AccommodationsHandler) GetAccommodationById(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	accommodationId := vars["id"]
	accommodations, err := a.repo.GetAccommodationById(accommodationId)
	if err != nil {
		a.logger.Print("Database exception: ", err)
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/accommodations", rw)
	}

	if accommodations == nil {
		return
	}

	if err := accommodations.ToJSON(rw); err != nil {
		a.handleJSONConversionError(errors.NewError("Conversion error", 500), rw)
		return
	}
}

func (a *AccommodationsHandler) handleJSONConversionError(err *errors.ErrorStruct, rw http.ResponseWriter) {
	http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
	a.logger.Fatal("Unable to convert to json:", err)
}

func (a *AccommodationsHandler) UpdateAccommodationById(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	accommodationId := vars["id"]
	location := vars["location"]
	UpdateAccommById := h.Context().Value(KeyProduct{}).(*domain.Accommodation)
	UpdateAccommById, err := a.repo.UpdateAccommodationById(accommodationId, location, UpdateAccommById)
	if err != nil {
		a.logger.Print("Database exception:", err)
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/accommodations", rw)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

func (a *AccommodationsHandler) DeleteAccommodationById(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	accommodationId := vars["id"]

	err, _ := a.repo.DeleteAccommodationById(accommodationId)
	if err != nil {
		a.logger.Print("Database exception: ", err)
		http.Error(rw, "Failed to delete the accommodation", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusTeapot)
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

func (a *AccommodationsHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		a.logger.Println("Method [", h.Method, "] - Hit path :", h.URL.Path)

		rw.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(rw, h)
	})
}
