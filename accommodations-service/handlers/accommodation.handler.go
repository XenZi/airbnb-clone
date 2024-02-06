package handlers

import (
	"accommodations-service/config"
	"accommodations-service/domain"
	"accommodations-service/services"
	"accommodations-service/utils"
	"encoding/csv"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
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
	Logger               *config.Logger
}

func NewAccommodationHandler(logger *config.Logger) *AccommodationsHandler {
	return &AccommodationsHandler{
		Logger: logger,
	}
}

func (a *AccommodationsHandler) CreateAccommodationById(rw http.ResponseWriter, h *http.Request) {
	//TODO tracing repair
	ctx, span := a.Tracer.Start(h.Context(), "AccommodationsHandler.CreateAccommodationById")
	defer span.End()
	//decoder.DisallowUnknownFields()
	var images []multipart.File

	contentType := h.Header.Get("Content-Type")
	isMultipart := strings.HasPrefix(contentType, "multipart/form-data")

	if isMultipart {
		err := h.ParseMultipartForm(50 << 20)
		if err != nil {
			a.Logger.Error("Error while parsing multipartForm", log.Fields{
				"module": "handler",
				"error":  err.Error(),
			})
			utils.WriteErrorResp(err.Error(), http.StatusBadRequest, "e puklo", rw)
			return
		}
		files, ok := h.MultipartForm.File["images"]
		if !ok || len(files) == 0 {
			a.Logger.Error("No files uploaded", log.Fields{
				"module": "handler",
				"error":  "No files Uploaded",
			})
			utils.WriteErrorResp(err.Error(), http.StatusBadRequest, "nista bajo", rw)
			return
		}
		for _, fileHeader := range files {
			file, err1 := fileHeader.Open()
			if err1 != nil {
				a.Logger.Error("Error returning formfile", log.Fields{
					"module": "handler",
					"error":  err1.Error(),
				})
				utils.WriteErrorResp(err.Error(), http.StatusBadRequest, "nista bajo", rw)
				return
			}
			defer file.Close()
			images = append(images, file)

		}
	}

	var accDates []domain.AvailableAccommodationDates
	datesJson := h.FormValue("availableAccommodationDates")
	err2 := json.Unmarshal([]byte(datesJson), &accDates)
	if err2 != nil {
		a.Logger.Error("Error unmarshaling json", log.Fields{
			"module": "handler",
			"error":  err2.Error(),
		})
		utils.WriteErrorResp(err2.Error(), http.StatusBadRequest, "dates puca", rw)
		return
	}
	minVis, err := strconv.Atoi(h.FormValue("minNumOfVisitors"))
	maxVis, err := strconv.Atoi(h.FormValue("maxNumOfVisitors"))

	var conv []string
	reader := csv.NewReader(strings.NewReader(h.FormValue("conveniences")))
	records, err := reader.ReadAll()
	if err != nil {
		a.Logger.Error("Error reading records", log.Fields{
			"module": "handler",
			"error":  err.Error(),
		})
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
		Paying:                      h.FormValue("paying"),
	}

	_, err4 := a.AccommodationService.CreateAccommodation(accomm, images, ctx)
	if err4 != nil {
		a.Logger.Error("Error creating accomodation", log.Fields{
			"module": "handler",
			"error":  err.Error(),
		})
		utils.WriteErrorResp(err4.GetErrorMessage(), 500, "ovo je druis", rw)
		return
	}
	a.Logger.Infof("Successfully sent accommodation to accommodation service")
	utils.WriteResp(accomm, 201, rw)
}

func (a *AccommodationsHandler) GetImage(rw http.ResponseWriter, r *http.Request) {
	ctx, span := a.Tracer.Start(r.Context(), "AccommodationsHandler.GetImage")
	defer span.End()
	vars := mux.Vars(r)
	imageId := vars["id"]
	file, err := a.AccommodationService.GetImage(ctx, imageId)
	if err != nil {
		a.Logger.Error("Error getting image", log.Fields{
			"module": "handler",
			"error":  err.GetErrorMessage(),
		})
		utils.WriteErrorResp(err.GetErrorMessage(), 500, "nemere slike otvarati", rw)
		a.Logger.Infof("Successfully got image")
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
		a.Logger.Error("Error getting accommodations", log.Fields{
			"module": "handler",
			"error":  err.GetErrorMessage(),
		})
		utils.WriteErrorResp(err.GetErrorMessage(), http.StatusInternalServerError, "api/accommodations", rw)
		return
	}
	a.Logger.Infof("Successfully got all accommodations")
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
		a.Logger.Error("Error getting accommodation by id", log.Fields{
			"module": "handler",
			"error":  err.GetErrorMessage(),
		})
		utils.WriteErrorResp(err.GetErrorMessage(), http.StatusNotFound, "api/accommodations/"+accommodationId, rw)
		return
	}

	// Serialize accommodation to JSON and write response

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	a.Logger.Infof("Successfully got accommodation by id" + accommodationId)

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
		a.Logger.Error("Error decoding ids", log.Fields{
			"module": "handler",
			"error":  decodeErr.Error(),
		})
		utils.WriteErrorResp(decodeErr.Error(), http.StatusInternalServerError, "api/accommodations/FindByIds", rw)
		return
	}
	log.Println("Idevi su", ids.Ids)
	accommodations, err := a.AccommodationService.FindAccommodationByIds(ctx, ids.Ids)
	if err != nil {
		a.Logger.Error("Error getting accommodation by list of ids", log.Fields{
			"module": "handler",
			"error":  err.GetErrorMessage(),
		})
		utils.WriteErrorResp(err.GetErrorMessage(), http.StatusNotFound, "api/accommodations/FindByIds", rw)
		return
	}

	// Serialize accommodation to JSON and write response
	a.Logger.Infof("Successfully got accommodation by list of ids")
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
		a.Logger.Error("Error getting accommodation by id in the update function", log.Fields{
			"module": "handler",
			"error":  err.GetErrorMessage(),
		})
		utils.WriteErrorResp(err.GetErrorMessage(), http.StatusInternalServerError, "api/accommodations/"+accommodationId, rw)
		return
	}

	var updatedAccommodation domain.Accommodation
	decodeErr := decoder.Decode(&updatedAccommodation)
	if decodeErr != nil {
		a.Logger.Error("Error decoding accommodation in the update function", log.Fields{
			"module": "handler",
			"error":  err.GetErrorMessage(),
		})
		utils.WriteErrorResp(decodeErr.Error(), http.StatusInternalServerError, "api/accommodations/"+accommodationId, rw)
		return
	}
	id, _ := primitive.ObjectIDFromHex(accommodationId)
	updatedAccommodation.Id = id

	accommodation, err := a.AccommodationService.UpdateAccommodation(ctx, updatedAccommodation)
	if err != nil {
		a.Logger.Error("Error getting response from accommodation service", log.Fields{
			"module": "handler",
			"error":  err.GetErrorMessage(),
		})
		utils.WriteErrorResp(err.GetErrorMessage(), http.StatusInternalServerError, "api/accommodations/"+accommodationId, rw)
		return
	}
	a.Logger.Infof("Successfully updated accommodation with the id" + accommodationId)
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
		a.Logger.Error("Error deleting accommodation with id:"+accommodationId, log.Fields{
			"module": "handler",
			"error":  err.GetErrorMessage(),
		})
		utils.WriteErrorResp(err.GetErrorMessage(), http.StatusInternalServerError, "api/accommodations/"+accommodationId, rw)
		return
	}
	a.Logger.Infof("Successfully deleted accommodation with id:" + accommodationId)
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
		a.Logger.Error("Error deleting accommodation with user id:"+userId, log.Fields{
			"module": "handler",
			"error":  err.GetErrorMessage(),
		})
		utils.WriteErrorResp(err.GetErrorMessage(), http.StatusInternalServerError, "api/accommodations/user"+userId, rw)
		return
	}
	a.Logger.Infof("Successfully deleted accommodation with user id:" + userId)
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
		a.Logger.Error("Error converting string into a number"+visitors, log.Fields{
			"module": "handler",
			"error":  err.Error(),
		})
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
		a.Logger.Error("Error searching accommodations", log.Fields{
			"module": "handler",
			"error":  errS.GetErrorMessage(),
		})
		utils.WriteErrorResp(errS.GetErrorMessage(), http.StatusInternalServerError, "api/accommodations/search", w)
		log.Println("greska je,", errS.GetErrorMessage())
		return
	}

	// Encode the search results into JSON and send the response
	//responseJSON, err := json.Marshal(accommodations)
	if err != nil {
		a.Logger.Error("Error marshaling JSON ", log.Fields{
			"module": "handler",
			"error":  err.Error(),
		})
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	a.Logger.Infof("Successfully passed the search function in handler")
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
		a.Logger.Error("Error decoding acccommodatin with id:"+accommodationID, log.Fields{
			"module": "handler",
			"error":  err.Error(),
		})
		log.Println(err)
		utils.WriteErrorResp("Internal server error", 500, "api/recommendation/rating/accommodation", w)
		return
	}
	log.Println(accommodation)
	// Now, you can use the 'rating' variable in your logic
	a.AccommodationService.PutAccommodationRating(ctx, accommodationID, accommodation)
	a.Logger.Infof("Successfully called PutAccommodationRating func for accommodation with id:" + accommodationID)
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
			a.Logger.Error("Error getting cache", log.Fields{
				"module": "handler",
				"error":  err.Error(),
			})
			next.ServeHTTP(rw, r)
		} else {
			a.Logger.Infof("Successfull cache hit")
			rw.Header().Set("Content-Type", "image/jpeg")
			rw.WriteHeader(http.StatusOK)
			rw.Write(image)
		}
	})
}
