package handlers

import (
	"github.com/gorilla/mux"
	"metrics_query/domain"
	"metrics_query/utils"
	"net/http"
)

type AccommodationHandler struct {
	store domain.AccommodationStore
}

func NewAccommodationHandler(store domain.AccommodationStore) AccommodationHandler {
	return AccommodationHandler{
		store: store,
	}
}

func (h AccommodationHandler) Get(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		utils.WriteErrorResp("bad request", 500, "metrics/get/{id}", rw)
		return
	}
	accommodation, err := h.store.Read(id)
	if err != nil {
		utils.WriteErrorResp(err.Error(), 404, "not found", rw)
		return
	}
	utils.WriteResp(accommodation, 200, rw)
	return
}

func (h AccommodationHandler) GetAll(rw http.ResponseWriter, r *http.Request) {
	accommodations := h.store.ReadAll()
	utils.WriteResp(accommodations, 200, rw)
}
