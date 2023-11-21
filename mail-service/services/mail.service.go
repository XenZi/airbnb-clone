package services

import (
	"github.com/XenZi/airbnb-clone/mail-service/domains"
	"github.com/XenZi/airbnb-clone/mail-service/errors"
)

type MailService struct {
	sender *domains.Sender
}

type ConfirmMailLink struct {
	Link string
}
const (
	CONFIRM_ACCOUNT_TEMPLATE string = "templates/confirm-account-template.html"
)

func NewMailService(sender *domains.Sender) *MailService {
	return &MailService{
		sender: sender,
	}
}

func (m MailService) SendRegisterConfirmationEmail(accountConfirmation domains.AccountConfirmation) *errors.ErrorStruct {
	if err := m.sender.SendHTMLEmail(
		CONFIRM_ACCOUNT_TEMPLATE,
		[]string{
			accountConfirmation.Email,
		},
		[]string{},
		"Account confirmation email",
		ConfirmMailLink{
			Link: "http://localhost:4200/confirm-account/" + accountConfirmation.Token,
		},
		[]string{}); err != nil {
		return err
	}
	return nil
}
