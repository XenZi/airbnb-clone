package utils

import (
	"encoding/json"
	"net/http"
	"time"
	"user-service/domain"
)

func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		WriteErrorResponse(err.Error(), http.StatusInternalServerError, "Internal Server Error", w)

	}
}

func WriteErrorResponse(err string, status int, path string, w http.ResponseWriter) {
	baseErrorResponse := domain.BaseErrorHttpResponse{
		Error:  err,
		Path:   path,
		Status: status,
		Time:   time.Now().String(),
	}
	writeJSONResponse(w, status, baseErrorResponse)
}

func WriteResp(resp any, statusCode int, w http.ResponseWriter) {
	if resp == nil {
		return
	}
	domainResponse := domain.BaseHttpResponse{
		Status: statusCode,
		Data:   resp,
	}
	writeJSONResponse(w, statusCode, domainResponse)
}
