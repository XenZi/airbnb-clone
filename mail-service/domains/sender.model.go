package domains

import (
	"bytes"
	"fmt"
	"github.com/XenZi/airbnb-clone/mail-service/errors"
	"html/template"
	"log"
	"mime/quotedprintable"
	"net/smtp"
	"strings"
	"time"
)

type Sender struct {
	Email          string
	Password       string
	SMTPServer     string
	SMTPServerPort int
}

func NewEmailSender(email, password, SMTPServer string, SMTPServerPort int) *Sender {
	return &Sender{Email: email, Password: password, SMTPServer: SMTPServer, SMTPServerPort: SMTPServerPort}
}

func (s *Sender) SendHTMLEmail(templatePath string, to, cc []string, subject string, data interface{}, files []string) *errors.ErrorStruct {
	if templatePath == "" {
		log.Println("Error while reading template, not exist on this path: " + templatePath)
		return errors.NewError("Error while reading template, not exist on this path", 500)
	}
	tmpl, err := ParseTemplate(templatePath, data)
	if err != nil {
		log.Println(err)
		return err
	}
	body := s.writeEmail(to, cc, "text/html", subject, tmpl, files)
	if err := s.sendEmail(to, subject, body); err != nil {
		log.Println(err)
		return errors.NewError(err.Error(), 500)
	}
	return nil
}
func (s *Sender) SendPlainEmail(to, cc []string, subject, data string, files []string) *errors.ErrorStruct {
	body := s.writeEmail(to, cc, "text/plain", subject, data, files)
	if err := s.sendEmail(to, subject, body); err != nil {
		log.Println(err)
		return errors.NewError(err.Error(), 500)
	}
	return nil
}

func (s *Sender) writeEmail(to, cc []string, ct, subj, body string, files []string) string {
	// Define variables.
	var message string
	var encodedBody bytes.Buffer

	// Create delimiter.
	delimiter := fmt.Sprintf("**=mail%d", time.Now().UnixNano())

	result := quotedprintable.NewWriter(&encodedBody) // Assuming encodedBody is a valid io.WriteCloser.

	if _, err := result.Write([]byte(body)); err != nil {
		log.Println(err)
	}
	defer func(result *quotedprintable.Writer) {
		err := result.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(result)

	// Create message.
	message += fmt.Sprintf("From: %s\r\n", s.Email)
	message += fmt.Sprintf("To: %s\r\n", strings.Join(to, ";"))
	if len(cc) > 0 {
		// If CC is specified.
		message += fmt.Sprintf("Cc: %s\r\n", strings.Join(cc, ";"))
	}
	message += fmt.Sprintf("Subject: %s\r\n", subj)
	message += fmt.Sprintf("MIME-Version: 1.0\r\n")
	message += fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", delimiter)

	// Add body of the message (from html template or plain text).
	message += fmt.Sprintf("--%s\r\n", delimiter)
	message += fmt.Sprintf("Content-Transfer-Encoding: quoted-printable\r\n")
	message += fmt.Sprintf("Content-Type: %s; charset=\"utf-8\"\r\n", ct)
	message += fmt.Sprintf("Content-Disposition: inline\r\n")
	message += fmt.Sprintf("%s\r\n", body)
	return message
}

func (s *Sender) sendEmail(dest []string, subject, body string) error {
	auth := smtp.PlainAuth("", s.Email, s.Password, s.SMTPServer)
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", s.SMTPServer, s.SMTPServerPort),
		auth,
		s.Email, dest, []byte(body),
	)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func ParseTemplate(file string, data interface{}) (string, *errors.ErrorStruct) {
	tmpl, errParseFiles := template.ParseFiles(file)
	if errParseFiles != nil {
		log.Println(errParseFiles)
		return "", errors.NewError(errParseFiles.Error(), 500)

	}
	buffer := new(bytes.Buffer)
	if errExecute := tmpl.Execute(buffer, data); errExecute != nil {
		log.Println(errExecute)
		return "", errors.NewError(errParseFiles.Error(), 500)
	}
	return buffer.String(), nil
}
