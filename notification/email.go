package notification

import (
	"net/smtp"
	"os"
)

type EmailMessage struct {
	Body     string
	Receiver []string
	Subject  string
	Mime     string
	Host     string
}

func AuthSmtpServer() smtp.Auth {
	auth := smtp.PlainAuth(
		"",
		os.Getenv("SENDER_EMAIL"),
		os.Getenv("SENDER_EMAIL_TOKEN"),
		os.Getenv("SENDER_EMAIL_HOST"),
	)
	return auth
}
func (e EmailMessage) SendEmail() error {
	message := []byte("subject:" + e.Subject + "\n" + e.Mime + e.Body)
	err := smtp.SendMail(
		e.Host, // is common to use gmai host if emails use send only gmail "smtp.gmail.com:587",
		AuthSmtpServer(),
		os.Getenv("SENDER_EMAIL"),
		e.Receiver,
		message,
	)
	return err
}
