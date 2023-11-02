package utils

import (
	"encoding/json"
	"env-test-app/domain"
	"env-test-app/errors"
	"net/http"
	"strings"
)

func WriteErrorResp(err error, w http.ResponseWriter) {
	if err == nil {
		return
	} else if err.Error() == errors.ErrUnathorized().Error() {
		w.WriteHeader(http.StatusForbidden)
	} else if strings.Contains(err.Error(), "not found") {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write([]byte(err.Error()))
}

func WriteResp(resp any, statusCode int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	if resp == nil {
		return
	}
	domainResponse := domain.BaseHttpResponse{Status: http.StatusOK, Data: resp}
	marshaledResponse, err := json.Marshal(domainResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(marshaledResponse)
}
