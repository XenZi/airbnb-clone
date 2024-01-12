package handler

import (
	"encoding/json"
	"net/http"
	"recommendation-service/domains"
	"recommendation-service/services"
	"recommendation-service/utils"

	"github.com/gorilla/mux"
)

type RatingHandler struct {
	service *services.RatingService
}

func NewRatingHandler(service *services.RatingService) *RatingHandler {
	return &RatingHandler{
		service: service,
	}
}

func (rh RatingHandler) CreateRatingForAccommodation(r http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var rating domains.RateAccommodation
	if err := decoder.Decode(&rating); err != nil {
		utils.WriteErrorResp("Internal server error", 500, "api/recommendation/rating/accommodation", r)
		return
	}
	ctx := h.Context()
	resp, err := rh.service.CreateRatingForAccommodation(ctx, rating)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/recommendation/rating/accommodation", r)
		return
	}
	utils.WriteResp(resp, 200, r)
}

func (rh RatingHandler) UpdateRatingForAccommodation(r http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var rating domains.RateAccommodation
	if err := decoder.Decode(&rating); err != nil {
		utils.WriteErrorResp("Internal server error", 500, "api/login", r)
	}
	ctx := h.Context()
	resp, err := rh.service.UpdateRatingForAccommodation(ctx, rating)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/recommendation/rating/accommodation", r)
		return
	}
	utils.WriteResp(resp, 200, r)
}

func (rh RatingHandler) DeleteRatingForAccommodation(r http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var rating domains.RateAccommodation
	if err := decoder.Decode(&rating); err != nil {
		utils.WriteErrorResp("Internal server error", 500, "api/login", r)
	}
	ctx := h.Context()
	resp, err := rh.service.DeleteRatingForAccommodation(ctx, rating)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/recommendation/rating/accommodation", r)
		return
	}
	utils.WriteResp(resp, 200, r)
}

func (rh RatingHandler) CreateRatingForHost(r http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var rating domains.RateHost
	if err := decoder.Decode(&rating); err != nil {
		utils.WriteErrorResp("Internal server error", 500, "api/login", r)
		return
	}
	resp, err := rh.service.CreateRatingForHost(rating)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/recommendations/ratings/host", r)
		return
	}
	utils.WriteResp(resp, 201, r)
}

func (rh RatingHandler) UpdateRatingForHost(r http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var rating domains.RateHost
	if err := decoder.Decode(&rating); err != nil {
		utils.WriteErrorResp("Internal server error", 500, "api/login", r)
		return
	}
	resp, err := rh.service.UpdateRatingForHostAndGuest(rating)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/recommendations/ratings/host", r)
		return
	}
	utils.WriteResp(resp, 201, r)
}

func (rh RatingHandler) DeleteRatingForHost(r http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var rating domains.RateHost
	if err := decoder.Decode(&rating); err != nil {
		utils.WriteErrorResp("Internal server error", 500, "api/login", r)
		return
	}
	resp, err := rh.service.DeleteRatingBetweenGuestAndHost(rating)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/recommendations/ratings/host", r)
		return
	}
	utils.WriteResp(resp, 201, r)
}

func (rh RatingHandler) GetAllRatingsForHost(r http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]
	if id == "" {
		utils.WriteErrorResp("Bad request", 400, "api/recommondations/host/"+id, r)
		return
	}
	resp, err := rh.service.GetAllRatingsForHostByID(id)
	if err != nil {
		utils.WriteErrorResp("Bad request", 400, "api/recommondations/host/"+id, r)
		return
	}
	utils.WriteResp(resp, 200, r)
}

func (rh RatingHandler) GetAllRatingsForAccommmodation(r http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]
	if id == "" {
		utils.WriteErrorResp("Bad request", 400, "api/recommondations/host/"+id, r)
		return
	}
	resp, err := rh.service.GetAllAccommodationRatings(id)
	if err != nil {
		utils.WriteErrorResp("Bad request", 400, "api/recommondations/host/"+id, r)
		return
	}
	utils.WriteResp(resp, 200, r)
}
