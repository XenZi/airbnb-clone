package handlers

import (
	"accommodations-service/domain"
	"accommodations-service/services"
	"accommodations-service/utils"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/trace"
)

type AccommodationsHandler struct {
	AccommodationService *services.AccommodationService
	Tracer               trace.Tracer
}

func (a *AccommodationsHandler) CreateAccommodationById(rw http.ResponseWriter, h *http.Request) {
	//TODO tracing reapair
	ctx, span := a.Tracer.Start(h.Context(), "AccommodationsHandler.CreateAccommodationById")
	defer span.End()
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

	_, err4 := a.AccommodationService.CreateAccommodation(accomm, image, ctx)
	if err4 != nil {
		utils.WriteErrorResp(err4.GetErrorMessage(), 500, "ovo je druis", rw)
		return
	}
	utils.WriteResp(accomm, 201, rw)
}

func (a *AccommodationsHandler) GetImage(rw http.ResponseWriter, r *http.Request) {
	ctx, span := a.Tracer.Start(r.Context(), "AccommodationsHandler.GetImage")
	defer span.End()
	vars := mux.Vars(r)
	imageId := vars["id"]
	file, err := a.AccommodationService.GetImage(ctx, imageId)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), 500, "nemere slike otvarati", rw)
		return
	}
	rw.Header().Set("Content-Type", "image/jpeg")
	rw.WriteHeader(http.StatusOK)
	rw.Write(file)

}

func (a *AccommodationsHandler) GetAllAccommodations(rw http.ResponseWriter, r *http.Request) {
	ctx, span := a.Tracer.Start(r.Context(), "AccommodationsHandler.GetAllAccommodations")
	defer span.End()
	accommodations, err := a.AccommodationService.GetAllAccommodations(ctx)
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
	ctx, span := a.Tracer.Start(r.Context(), "AccommodationsHandler.GetAccommodationById")
	defer span.End()
	vars := mux.Vars(r)
	accommodationId := vars["id"]

	accommodation, err := a.AccommodationService.GetAccommodationById(ctx, accommodationId)
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
	ctx, span := a.Tracer.Start(r.Context(), "AccommodationsHandler.FindAccommodationsByIds")
	defer span.End()
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
	accommodations, err := a.AccommodationService.FindAccommodationByIds(ctx, ids.Ids)
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
	ctx, span := a.Tracer.Start(r.Context(), "AccommodationsHandler.UpdateAccommodationById")
	defer span.End()
	vars := mux.Vars(r)
	accommodationId := vars["id"]
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	_, err := a.AccommodationService.GetAccommodationById(ctx, accommodationId)
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

	accommodation, err := a.AccommodationService.UpdateAccommodation(ctx, updatedAccommodation)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), http.StatusInternalServerError, "api/accommodations/"+accommodationId, rw)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	utils.WriteResp(accommodation, 201, rw)
}

func (a *AccommodationsHandler) DeleteAccommodationById(rw http.ResponseWriter, r *http.Request) {
	ctx, span := a.Tracer.Start(r.Context(), "AccommodationsHandler.DeleteAccommodationById")
	defer span.End()
	vars := mux.Vars(r)
	accommodationId := vars["id"]

	_, err := a.AccommodationService.DeleteAccommodation(ctx, accommodationId)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), http.StatusInternalServerError, "api/accommodations/"+accommodationId, rw)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusNoContent) // HTTP 204 No Content for successful deletion
}

func (a *AccommodationsHandler) DeleteAccommodationsByUserId(rw http.ResponseWriter, r *http.Request) {
	ctx, span := a.Tracer.Start(r.Context(), "AccommodationsHandler.DeleteAccommodationsByUserId")
	defer span.End()
	vars := mux.Vars(r)
	userId := vars["id"]
	log.Println("user id je:", userId)

	err := a.AccommodationService.DeleteAccommodationsByUserId(ctx, userId)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), http.StatusInternalServerError, "api/accommodations/user"+userId, rw)
		return
	}

	utils.WriteResp("successfully deleted accommodations", 201, rw)
}

func (a *AccommodationsHandler) SearchAccommodations(w http.ResponseWriter, r *http.Request) {
	ctx, span := a.Tracer.Start(r.Context(), "AccommodationsHandler.SearchAccommodations")
	defer span.End()

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

	accommodations, errS := a.AccommodationService.SearchAccommodations(city, country, numOfVisitors, startDate, endDate, maxPrice, conveniences, isDistinguishedString, ctx)

	if errS != nil {
		utils.WriteErrorResp(errS.GetErrorMessage(), http.StatusInternalServerError, "api/accommodations/search", w)
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
	ctx, span := a.Tracer.Start(r.Context(), "AccommodationsHandler.PutAccommodationRating")
	defer span.End()
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
	a.AccommodationService.PutAccommodationRating(ctx, accommodationID, accommodation)

	// Respond with a success message or any appropriate response
	w.Header().Set("Content-Type", "application/json")
	utils.WriteResp(accommodation, 201, w)
}

func (a *AccommodationsHandler) MiddlewareCacheHit(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx, span := a.Tracer.Start(r.Context(), "AccommodationsHandler.MiddlewareCacheHit")
		defer span.End()
		vars := mux.Vars(r)
		id := vars["id"]
		image, err := a.AccommodationService.GetCache(ctx, id)
		if err != nil {
			next.ServeHTTP(rw, r)
		} else {
			rw.Header().Set("Content-Type", "image/jpeg")
			rw.WriteHeader(http.StatusOK)
			rw.Write(image)
		}
	})
}
