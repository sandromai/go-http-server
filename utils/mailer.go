package utils

import (
	"net/smtp"
	"strings"

	"github.com/sandromai/go-http-server/types"
)

type Mailer struct {
	Host     string
	Port     string
	Username string
	Password string
}

func (mailer *Mailer) Send(
	fromEmail,
	fromName,
	toEmail,
	toName,
	subject,
	body string,
) *types.AppError {
	if !CheckEmail(fromEmail) {
		return &types.AppError{
			StatusCode: 500,
			Message:    "Invalid origin email address.",
		}
	}

	if !CheckEmail(toEmail) {
		return &types.AppError{
			StatusCode: 500,
			Message:    "Invalid destination email address.",
		}
	}

	formattedMessage := mailer.formatBody(fromEmail,
		fromName,
		toEmail,
		toName,
		subject,
		body,
	)

	auth := smtp.PlainAuth("", mailer.Username, mailer.Password, mailer.Host)

	if err := smtp.SendMail(mailer.Host+":"+mailer.Port, auth, mailer.Username, []string{toEmail}, formattedMessage); err != nil {
		return &types.AppError{
			StatusCode: 500,
			Message:    "Error sending email.",
		}
	}

	return nil
}

func (*Mailer) formatBody(
	fromEmail,
	fromName,
	toEmail,
	toName,
	subject,
	body string,
) (formattedMessage []byte) {
	if fromName == "" {
		formattedMessage = append(formattedMessage, []byte("From: "+fromEmail+"\n")...)
	} else {
		formattedMessage = append(formattedMessage, []byte("From: "+strings.TrimSpace(fromName)+" <"+fromEmail+">\n")...)
	}

	if toName == "" {
		formattedMessage = append(formattedMessage, []byte("To: "+toEmail+"\n")...)
	} else {
		formattedMessage = append(formattedMessage, []byte("To: "+strings.TrimSpace(toName)+" <"+toEmail+">\n")...)
	}

	formattedMessage = append(formattedMessage, []byte("Subject: "+strings.TrimSpace(subject)+"\n")...)
	formattedMessage = append(formattedMessage, []byte("MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n")...)
	formattedMessage = append(formattedMessage, []byte(strings.TrimSpace(body))...)

	return formattedMessage
}
