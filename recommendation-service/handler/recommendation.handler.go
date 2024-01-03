package handler

import (
	"net/http"
	"recommendation-service/services"
	"recommendation-service/utils"

	"github.com/gorilla/mux"
)

type RecommendationHandler struct {
	service *services.RecommendationService
}

func NewRecommendationHandler(service *services.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{
		service: service,
	}
}

func (rh RecommendationHandler) GetAllRecommendationsForUser(r http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]
	if id == "" {
		utils.WriteErrorResp("Bad request", 400, "api/recommondations/host/"+id, r)
		return
	}
	recommendations, err := rh.service.GetAllRecommendationsByUserID(id)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/recommendations", r)
		return
	}
	utils.WriteResp(recommendations, 200, r)
}
