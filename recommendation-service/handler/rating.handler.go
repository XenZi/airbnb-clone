package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"recommendation-service/domains"
	"recommendation-service/services"
	"recommendation-service/utils"
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
		utils.WriteErrorResp("Internal server error", 500, "api/login", r)
	}
	rh.service.CreateRatingForAccommodation(rating)
}

func (rh RatingHandler) CreateRatingForHost(r http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var rating domains.RateHost
	if err := decoder.Decode(&rating); err != nil {
		utils.WriteErrorResp("Internal server error", 500, "api/login", r)
		return;
	}
	resp, err := rh.service.CreateRatingForHost(rating)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/recommendations/ratings/host", r)
		return
	}
	utils.WriteResp(resp, 201, r)
}

func (rh RatingHandler) UpdateRatingForHost(r http.ResponseWriter, h *http.Request) {
	log.Println("SADSADASDASD")
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var rating domains.RateHost
	if err := decoder.Decode(&rating); err != nil {
		utils.WriteErrorResp("Internal server error", 500, "api/login", r)
		return;
	}
	resp, err := rh.service.UpdateRatingForHostAndGuest(rating)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/recommendations/ratings/host", r)
		return
	}
	utils.WriteResp(resp, 201, r)
}