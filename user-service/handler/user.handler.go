package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"user-service/config"
	"user-service/domain"
	"user-service/service"
	"user-service/utils"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	UserService *service.UserService
	logger      *config.Logger
}

const source = "user-handler"

func NewUserHandler(userService *service.UserService, logger *config.Logger) *UserHandler {
	return &UserHandler{
		UserService: userService,
		logger:      logger,
	}
}

func (u UserHandler) CreateHandler(rw http.ResponseWriter, h *http.Request) {
	u.logger.LogInfo(source, "Received Create request")
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var createData domain.CreateUser
	if err := decoder.Decode(&createData); err != nil {
		utils.WriteErrorResponse(err.Error(), 500, "api/create", rw)
		return
	}
	user, err := u.UserService.CreateUser(createData)
	if err != nil {
		utils.WriteErrorResponse(err.GetErrorMessage(), err.GetErrorStatus(), "api/create", rw)
		return
	}
	utils.WriteResp(user, 200, rw)
}

func (u UserHandler) UpdateHandler(rw http.ResponseWriter, h *http.Request) {
	u.logger.LogInfo(source, "Received Update request")
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var updateData domain.CreateUser
	if err := decoder.Decode(&updateData); err != nil {
		utils.WriteErrorResponse(err.Error(), 500, "api/update", rw)
		return
	}
	user, err := u.UserService.UpdateUser(updateData)
	if err != nil {
		utils.WriteErrorResponse(err.GetErrorMessage(), err.GetErrorStatus(), "api/update", rw)
		return
	}
	utils.WriteResp(user, 200, rw)
}

func (u UserHandler) CredsHandler(rw http.ResponseWriter, h *http.Request) {
	u.logger.LogInfo(source, "Received Update Credentials request")
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var updateData domain.CreateUser
	if err := decoder.Decode(&updateData); err != nil {
		utils.WriteErrorResponse(err.Error(), 500, "api/update", rw)
		return
	}
	user, err := u.UserService.UpdateUserCreds(updateData)
	if err != nil {
		utils.WriteErrorResponse(err.GetErrorMessage(), err.GetErrorStatus(), "api/update", rw)
		return
	}
	utils.WriteResp(user, 200, rw)
}

func (u UserHandler) GetAllHandler(rw http.ResponseWriter, h *http.Request) {
	u.logger.LogInfo(source, "Received Get All request")
	userCollection, err := u.UserService.GetAllUsers()
	if err != nil {
		utils.WriteErrorResponse(err.GetErrorMessage(), err.GetErrorStatus(), "api/get-all", rw)
		return
	}
	utils.WriteResp(userCollection, 200, rw)
}

func (u UserHandler) GetUserById(rw http.ResponseWriter, h *http.Request) {
	u.logger.LogInfo(source, "Received Get By ID request")
	vars := mux.Vars(h)
	id := vars["id"]
	log.Println("Id koji preuzimam iz urla je,", id)
	user, hostUser, err := u.UserService.GetUserById(id)
	if err != nil {
		utils.WriteErrorResponse(err.GetErrorMessage(), err.GetErrorStatus(), "api/get-user", rw)
		return
	}
	if hostUser != nil {
		utils.WriteResp(hostUser, 200, rw)
		return
	}
	utils.WriteResp(user, 200, rw)
}

func (u UserHandler) UpdateRating(rw http.ResponseWriter, h *http.Request) {
	u.logger.LogInfo(source, "Received Update Rating request")
	vars := mux.Vars(h)
	id := vars["id"]
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var rating *domain.RatingStruct
	if err := decoder.Decode(&rating); err != nil {
		utils.WriteErrorResponse(err.Error(), 500, "api/users/rating", rw)
		return
	}
	log.Println(rating)
	ratingStr := rating.Rating
	log.Println(ratingStr)
	ratingF, err := strconv.ParseFloat(ratingStr, 64)
	if err != nil {
		utils.WriteErrorResponse("cannot convert to float64", 400, "api/users/rating", rw)
		return
	}
	erro := u.UserService.UpdateRating(id, ratingF)
	if erro != nil {
		utils.WriteErrorResponse(erro.GetErrorMessage(), erro.GetErrorStatus(), "api/users/rating", rw)
		return
	}
	utils.WriteResp(id, 200, rw)
}

func (u UserHandler) DeleteHandler(rw http.ResponseWriter, h *http.Request) {
	u.logger.LogInfo(source, "Received Delete request")
	vars := mux.Vars(h)
	id := vars["id"]
	ctx := h.Context()
	role := ctx.Value("role")
	log.Println("DELETED USER ROLE: ", role.(string))
	err := u.UserService.DeleteUser(role.(string), id)
	if err != nil {
		utils.WriteErrorResponse(err.GetErrorMessage(), err.GetErrorStatus(), "api/delete", rw)
		return
	}
	utils.WriteResp(id, 200, rw)
}
