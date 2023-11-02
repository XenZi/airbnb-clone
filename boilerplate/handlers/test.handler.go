package handlers

import (
	"env-test-app/services"
	"env-test-app/utils"
	"net/http"
)

type TestHandler struct {
	Service services.TestService
}

func (t TestHandler) SayHiFromHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteResp(t.Service.CallService(), 200, w)
}
