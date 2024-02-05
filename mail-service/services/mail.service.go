package services

import (
	"fmt"

	"github.com/XenZi/airbnb-clone/mail-service/config"
	"github.com/XenZi/airbnb-clone/mail-service/domains"
	"github.com/XenZi/airbnb-clone/mail-service/errors"
	"go.opentelemetry.io/otel/trace"
)

type MailService struct {
	sender  *domains.Sender
	logger  *config.Logger
	tracing trace.Tracer
}

type ConfirmMailLink struct {
	Link string
}

const (
	CONFIRM_ACCOUNT_TEMPLATE       string = "templates/confirm-account-template.html"
	RESET_PASSWORD_TEMPLATE        string = "templates/reset-password-template.html"
	NOTIFICATION_PROVIDER_TEMPLATE string = "templates/notification-provider-template.html"
)

func NewMailService(sender *domains.Sender, logger *config.Logger, tracing trace.Tracer) *MailService {
	return &MailService{
		sender:  sender,
		logger:  logger,
		tracing: tracing,
	}
}

func (m MailService) SendRegisterConfirmationEmail(accountConfirmation domains.AccountConfirmation) *errors.ErrorStruct {
	m.logger.LogInfo("mail-service", fmt.Sprintf("Sending user confirmation mail to %v", accountConfirmation.Email))
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
		m.logger.LogError("mail-service", fmt.Sprintf("Error while sending email to %v", accountConfirmation.Email))
		return err
	}
	m.logger.LogInfo("mail-service", fmt.Sprintf("Confirmation mail successfully sent to %v", accountConfirmation.Email))

	return nil
}

func (m MailService) SendPasswordReset(requestResetPassword domains.RequestResetPassword) *errors.ErrorStruct {
	m.logger.LogInfo("mail-service", fmt.Sprintf("Sending user password reset mail to %v", requestResetPassword.Email))
	if err := m.sender.SendHTMLEmail(
		RESET_PASSWORD_TEMPLATE,
		[]string{
			requestResetPassword.Email,
		},
		[]string{},
		"Password reset mail",
		ConfirmMailLink{
			Link: "http://localhost:4200/reset-password/" + requestResetPassword.Token,
		},
		[]string{}); err != nil {
		m.logger.LogError("mail-service", fmt.Sprintf("Error while sending email to %v", requestResetPassword.Email))

		return err
	}
	m.logger.LogInfo("mail-service", fmt.Sprintf("Successfully user password reset mail to %v", requestResetPassword.Email))
	return nil
}

func (m MailService) SendNotification(mailNotification domains.NotificationMail) *errors.ErrorStruct {
	m.logger.LogInfo("mail-service", fmt.Sprintf("Sending user notification mail to %v", mailNotification.Email))
	if err := m.sender.SendHTMLEmail(
		NOTIFICATION_PROVIDER_TEMPLATE,
		[]string{
			mailNotification.Email,
		},
		[]string{},
		"Notification",
		mailNotification,
		[]string{}); err != nil {
		m.logger.LogError("mail-service", fmt.Sprintf("Error while sending email to %v", mailNotification.Email))
		return err
	}
	m.logger.LogInfo("mail-service", fmt.Sprintf("Successfully notification mail sent to %v", mailNotification.Email))
	return nil
}
