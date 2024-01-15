package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"user-service/domain"
	"user-service/service"
	"user-service/utils"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	UserService *service.UserService
}

func (u UserHandler) CreateHandler(rw http.ResponseWriter, h *http.Request) {
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
	userCollection, err := u.UserService.GetAllUsers()
	if err != nil {
		utils.WriteErrorResponse(err.GetErrorMessage(), err.GetErrorStatus(), "api/get-all", rw)
		return
	}
	utils.WriteResp(userCollection, 200, rw)
}

func (u UserHandler) GetUserById(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]
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
	vars := mux.Vars(h)
	id := vars["id"]
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var rating *domain.RatingStruct
	if err := decoder.Decode(&rating); err != nil {
		utils.WriteErrorResponse(err.Error(), 500, "api/users/rating", rw)
		return
	}
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
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusNoContent)
}

func (u UserHandler) DeleteHandler(rw http.ResponseWriter, h *http.Request) {
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
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusNoContent)
}

func (p *ProductsHandler) MiddlewareCacheHit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		vars := mux.Vars(h)
		id := vars["id"]
		// NoSQL: first look in the cache
		product, err := p.cache.Get(id)
		if err != nil {
			// If Product not present in cache, continue execution to handler method
			next.ServeHTTP(rw, h)
		} else {
			// If Product present in cache, return Product from cache
			err = product.ToJSON(rw)
			if err != nil {
				http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
				p.logger.Fatal("Unable to convert to json :", err)
				return
			}
		}
	})
}
