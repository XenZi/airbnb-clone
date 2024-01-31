package handlers

import (
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
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
	period, ok := vars["period"]
	if !ok {
		utils.WriteErrorResp("bad request", 500, "metrics/get/{id}", rw)
		return
	}
	accommodation, err := h.store.Read(id, period)
	if err != nil {
		utils.WriteErrorResp(err.Error(), 404, "not found", rw)
		return
	}
	utils.WriteResp(accommodation, 200, rw)
	return
}

func (h AccommodationHandler) GenUUID(rw http.ResponseWriter, r *http.Request) {
	dat := primitive.NewObjectID()
	hex := dat.Hex()
	log.Println(dat)
	log.Println(hex)
	dathex, _ := primitive.ObjectIDFromHex(hex)

	log.Println(dathex)

	utils.WriteResp(dat, 200, rw)
}
