package handlers

import (
	"accommodations-service/domain"
	"accommodations-service/services"
	"accommodations-service/utils"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccommodationsHandler struct {
	AccommodationService *services.AccommodationService
}

func (a *AccommodationsHandler) CreateAccommodationById(rw http.ResponseWriter, h *http.Request) {
	//decoder.DisallowUnknownFields()
	var image multipart.File

	contentType := h.Header.Get("Content-Type")
	isMultipart := strings.HasPrefix(contentType, "multipart/form-data")

	if isMultipart {
		err := h.ParseMultipartForm(10 << 20)
		if err != nil {
			utils.WriteErrorResp(err.Error(), http.StatusBadRequest, "e puklo", rw)
			return
		}
		file, _, err := h.FormFile("images")
		if err != nil {
			utils.WriteErrorResp(err.Error(), http.StatusBadRequest, "nista bajo", rw)
			return
		}
		image = file
		defer file.Close()
	}

	var accDates []domain.AvailableAccommodationDates
	datesJson := h.FormValue("availableAccommodationDates")
	err2 := json.Unmarshal([]byte(datesJson), &accDates)
	if err2 != nil {
		utils.WriteErrorResp(err2.Error(), http.StatusBadRequest, "dates puca", rw)
		return
	}
	minVis, err := strconv.Atoi(h.FormValue("minNumOfVisitors"))
	maxVis, err := strconv.Atoi(h.FormValue("maxNumOfVisitors"))

	var conv []string
	reader := csv.NewReader(strings.NewReader(h.FormValue("conveniences")))
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}
	for _, row := range records {
		for _, value := range row {
			conv = append(conv, value)
		}
	}
	accomm := domain.CreateAccommodation{
		Name:                        h.FormValue("name"),
		Address:                     h.FormValue("address"),
		City:                        h.FormValue("city"),
		Country:                     h.FormValue("country"),
		UserName:                    h.FormValue("username"),
		UserId:                      h.FormValue("userId"),
		Email:                       h.FormValue("email"),
		Conveniences:                conv,
		MinNumOfVisitors:            minVis,
		MaxNumOfVisitors:            maxVis,
		AvailableAccommodationDates: accDates,
		Location:                    h.FormValue("location"),
	}

	ctx, cancel := context.WithTimeout(h.Context(), time.Second*5)
	defer cancel()
	_, err4 := a.AccommodationService.CreateAccommodation(accomm, image, ctx)
	if err4 != nil {
		utils.WriteErrorResp(err4.GetErrorMessage(), 500, "ovo je druis", rw)
		return
	}
	utils.WriteResp(accomm, 201, rw)
}

func (a *AccommodationsHandler) GetImage(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageId := vars["id"]
	file, err := a.AccommodationService.GetImage(imageId)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), 500, "nemere slike otvarati", rw)
		return
	}
	rw.Header().Set("Content-Type", "image/jpeg")
	rw.WriteHeader(http.StatusOK)
	rw.Write(file)

}

func (a *AccommodationsHandler) GetAllAccommodations(rw http.ResponseWriter, r *http.Request) {
	accommodations, err := a.AccommodationService.GetAllAccommodations()
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), http.StatusInternalServerError, "api/accommodations", rw)
		return
	}
	//Serialize accommodations to JSON and write response
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

func (a *AccommodationsHandler) FindAccommodationsByIds(rw http.ResponseWriter, r *http.Request) {
	log.Println("Da li ulazi?")
	decoder := json.NewDecoder(r.Body)

	type IdStruct struct {
		Ids []string `json:"ids"`
	}
	var ids IdStruct
	decodeErr := decoder.Decode(&ids)
	if decodeErr != nil {
		utils.WriteErrorResp(decodeErr.Error(), http.StatusInternalServerError, "api/accommodations/FindByIds", rw)
		return
	}
	log.Println("Idevi su", ids.Ids)
	accommodations, err := a.AccommodationService.FindAccommodationByIds(ids.Ids)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), http.StatusNotFound, "api/accommodations/FindByIds", rw)
		return
	}

	// Serialize accommodation to JSON and write response

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	utils.WriteResp(accommodations, 201, rw)
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

func (a *AccommodationsHandler) DeleteAccommodationsByUserId(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]
	log.Println("user id je:", userId)

	err := a.AccommodationService.DeleteAccommodationsByUserId(userId)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), http.StatusInternalServerError, "api/accommodations/user"+userId, rw)
		return
	}

	utils.WriteResp("successfully deleted accommodations", 201, rw)
}

func (a *AccommodationsHandler) SearchAccommodations(w http.ResponseWriter, r *http.Request) {

	city := r.URL.Query().Get("city")
	log.Println("grad je", city)
	country := r.URL.Query().Get("country")

	visitors := r.URL.Query().Get("numOfVisitors")
	if visitors == "" {
		visitors = "1"
	}
	numOfVisitors, err := strconv.Atoi(visitors)

	if err != nil {
		utils.WriteErrorResp(err.Error(), 500, "api/accommodations/search", w)
		return
	}

	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")
	log.Println(startDate, endDate)

	maxPriceString := r.URL.Query().Get("maxPrice")

	if maxPriceString == "" {
		maxPriceString = "0"
	}

	maxPrice, err := strconv.Atoi(maxPriceString)

	conveniencesCsv := r.URL.Query().Get("conveniences")

	var conveniences []string
	if conveniencesCsv != "" {
		conveniences = strings.Split(conveniencesCsv, ",")
	}

	isDistinguishedString := r.URL.Query().Get("isDistinguished")
	log.Println("Is distinguished string", isDistinguishedString)

	// Handle empty dateRange as needed

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()
	accommodations, errS := a.AccommodationService.SearchAccommodations(city, country, numOfVisitors, startDate, endDate, maxPrice, conveniences, isDistinguishedString, ctx)

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

func (a *AccommodationsHandler) PutAccommodationRating(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accommodationID := vars["id"]

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var accommodation domain.Accommodation
	log.Println("SADSADDSADAS")
	if err := decoder.Decode(&accommodation); err != nil {
		log.Println(err)
		utils.WriteErrorResp("Internal server error", 500, "api/recommendation/rating/accommodation", w)
		return
	}
	log.Println(accommodation)
	// Now, you can use the 'rating' variable in your logic
	a.AccommodationService.PutAccommodationRating(accommodationID, accommodation)

	// Respond with a success message or any appropriate response
	w.Header().Set("Content-Type", "application/json")
	utils.WriteResp(accommodation, 201, w)
}

func (a *AccommodationsHandler) MiddlewareCacheHit(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		image, err := a.AccommodationService.GetCache(id)
		if err != nil {
			next.ServeHTTP(rw, r)
		} else {
			rw.Header().Set("Content-Type", "image/jpeg")
			rw.WriteHeader(http.StatusOK)
			rw.Write(image)
		}
	})
}
