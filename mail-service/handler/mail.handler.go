package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/XenZi/airbnb-clone/mail-service/domains"
	"github.com/XenZi/airbnb-clone/mail-service/services"
	"github.com/XenZi/airbnb-clone/mail-service/utils"
)

type MailHandler struct {
	mailService *services.MailService
}

func NewMailHandler(mailService *services.MailService) *MailHandler {
	return &MailHandler{
		mailService: mailService,
	}
}

func (m MailHandler) SendAccountConfirmationEmail(rw http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var accountConfirmationData domains.AccountConfirmation
	if err := decoder.Decode(&accountConfirmationData); err != nil {
		utils.WriteErrorResp(err.Error(), 500, "api/login", rw)
		return
	}
	log.Println(accountConfirmationData)
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

func (m MailHandler) SendPasswordResetEmail(rw http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var requestResetPassword domains.RequestResetPassword
	if err := decoder.Decode(&requestResetPassword); err != nil {
		utils.WriteErrorResp(err.Error(), 500, "api/login", rw)
		return
	}
	err := m.mailService.SendPasswordReset(requestResetPassword)
	if err != nil {
		log.Println("ERRR")
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "api/mail", rw)
		return
	}
}

func (m MailHandler) SendNotification(rw http.ResponseWriter, h *http.Request) {
	decoder := json.NewDecoder(h.Body)
	decoder.DisallowUnknownFields()
	var mailNotification domains.NotificationMail
	if err := decoder.Decode(&mailNotification); err != nil {
		utils.WriteErrorResp(err.Error(), 500, "api/login", rw)
		return
	}
	err := m.mailService.SendNotification(mailNotification)
	if err != nil {
		utils.WriteErrorResp(err.GetErrorMessage(), err.GetErrorStatus(), "/api/send-notification", rw)
		return
	}
	utils.WriteResp("Successfully gone", 200, rw)
}