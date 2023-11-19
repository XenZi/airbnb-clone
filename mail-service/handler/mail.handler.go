package handler

import (
	"encoding/json"
	"github.com/XenZi/airbnb-clone/mail-service/domains"
	"github.com/XenZi/airbnb-clone/mail-service/services"
	"github.com/XenZi/airbnb-clone/mail-service/utils"
	"log"
	"net/http"
)

type MailHandler struct {
	mailService *services.MailService
}

func NewMailHandler(mailService *services.MailService) *MailHandler {
	return &MailHandler{
		mailService: mailService,
	}
}

func (m MailHandler) SubmitAccountConfirmationMail(rw http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var accountConfirmationData domains.AccountConfirmation
	if err := decoder.Decode(&accountConfirmationData); err != nil {
		utils.WriteErrorResp(err.Error(), 500, "api/login", rw)
		return
	}
	err := m.mailService.SendRegisterConfirmationEmail(accountConfirmationData)
	if err != nil {
		log.Println("ERRR")
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/mail", rw)
		return
	}
	utils.WriteResp(map[string]string{
		"test": "testadsjkoadskjasdjlk",
	}, 200, rw)
	log.Println("SADDSADSASDAADSDS")
}
