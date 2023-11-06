package utils

import (
	"auth-service/domains"
	"encoding/json"
	"net/http"
	"time"
)

func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Handle JSON encoding error
		WriteErrorResp(err, http.StatusInternalServerError, "Internal Server Error", w)
	}
}

func WriteErrorResp(err error, status int, path string, w http.ResponseWriter) {
	if err == nil {
		return
	}
	baseErrorResp := domains.BaseErrorHttpResponse{
		Error:  err.Error(),
		Path:   path,
		Status: status,
		Time:   time.Now().String(),
	}
	writeJSONResponse(w, status, baseErrorResp)
}

func WriteResp(resp any, statusCode int, w http.ResponseWriter) {
	if resp == nil {
		return
	}
	domainResponse := domains.BaseHttpResponse{
		Status: statusCode,
		Data:   resp,
	}
	writeJSONResponse(w, statusCode, domainResponse)
}
