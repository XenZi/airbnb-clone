package handlers

import (
	"accommodations-service/domain"
	"accommodations-service/services"
	"accommodations-service/utils"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccommodationsHandler struct {
	AccommodationService *services.AccommodationService
}

func (a *AccommodationsHandler) CreateAccommodationById(rw http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var accomm domain.CreateAccommodation
	if err := decoder.Decode(&accomm); err != nil {
		utils.WriteErrorResp(err.Error(), 500, "api/accommodations", rw)
		return
	}
	ctx, cancel := context.WithTimeout(h.Context(), time.Second*5)
	defer cancel()
	accommodation, err := a.AccommodationService.CreateAccommodation(accomm, ctx)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), 500, "api/accommodations", rw)
	}
	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(http.StatusOK)
	utils.WriteResp(accommodation, 201, rw)

}

func (a *AccommodationsHandler) GetAllAccommodations(rw http.ResponseWriter, r *http.Request) {
	accommodations, err := a.AccommodationService.GetAllAccommodations()
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), http.StatusInternalServerError, "api/accommodations", rw)
		return
	}

	// Serialize accommodations to JSON and write response

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	utils.WriteResp(accommodations, 201, rw)
}

func (a *AccommodationsHandler) GetAccommodationById(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accommodationId := vars["id"]

	accommodation, err := a.AccommodationService.GetAccommodationById(accommodationId)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), http.StatusNotFound, "api/accommodations/"+accommodationId, rw)
		return
	}

	// Serialize accommodation to JSON and write response

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	utils.WriteResp(accommodation, 201, rw)
}

func (a *AccommodationsHandler) UpdateAccommodationById(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accommodationId := vars["id"]
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	_, err := a.AccommodationService.GetAccommodationById(accommodationId)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), http.StatusInternalServerError, "api/accommodations/"+accommodationId, rw)
		return
	}

	var updatedAccommodation domain.Accommodation
	decodeErr := decoder.Decode(&updatedAccommodation)
	if decodeErr != nil {
		utils.WriteErrorResp(decodeErr.Error(), http.StatusInternalServerError, "api/accommodations/"+accommodationId, rw)
		return
	}
	id, _ := primitive.ObjectIDFromHex(accommodationId)
	updatedAccommodation.Id = id

	accommodation, err := a.AccommodationService.UpdateAccommodation(updatedAccommodation)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), http.StatusInternalServerError, "api/accommodations/"+accommodationId, rw)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	utils.WriteResp(accommodation, 201, rw)
}

func (a *AccommodationsHandler) DeleteAccommodationById(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accommodationId := vars["id"]

	_, err := a.AccommodationService.DeleteAccommodation(accommodationId)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), http.StatusInternalServerError, "api/accommodations/"+accommodationId, rw)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusNoContent) // HTTP 204 No Content for successful deletion
}

func (a *AccommodationsHandler) SearchAccommodations(w http.ResponseWriter, r *http.Request) {

	city := r.URL.Query().Get("city")
	log.Println("grad je", city)
	country := r.URL.Query().Get("country")
	address := r.URL.Query().Get("address")
	visitors := r.URL.Query().Get("numOfVisitors")
	if visitors == "" {
		visitors = "1"
	}
	numOfVisitors, err := strconv.Atoi(visitors)

	if err != nil {
		utils.WriteErrorResp(err.Error(), 500, "api/accommodations/search", w)
		return
	}

	// Call the AccommodationService to perform the search
	accommodations, errS := a.AccommodationService.SearchAccommodations(city, country, address, numOfVisitors)

	if errS != nil {
		utils.WriteErrorResp(errS.GetErrorMessage(), http.StatusInternalServerError, "api/accommodations/BILOSTA", w)
		log.Println("greska je,", errS.GetErrorMessage())
		return
	}

	// Encode the search results into JSON and send the response
	//responseJSON, err := json.Marshal(accommodations)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	utils.WriteResp(accommodations, 201, w)

}
