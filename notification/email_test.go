package notification

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/mohamedabdifitah/processor/template"
)

func TestSendEmail(t *testing.T) {
	err := godotenv.Load("../local.env")
	if err != nil {
		log.Fatal(err)
	}
	e := EmailMessage{
		Body:     "Helllo world\n",
		Receiver: []string{os.Getenv("SENDER_EMAIL")},
		Subject:  "Testing Emails",
		Host:     "smtp.gmail.com:587",
	}
	err = e.SendEmail()
	if err != nil {
		t.Error(err)
	}
	template := template.AllTemplates
	template.LoadTemplates("../template/template.json")

	message, err := template.TempelateInjector(
		"OtpTemplate",
		map[string]string{
			"ExpireTime": "30",
			"Otp":        "55567",
			"Unit":       "minutes",
		},
	)
	if err != nil {
		t.Error(err)
	}
	e2 := EmailMessage{
		Body:     message,
		Mime:     "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\"",
		Receiver: []string{os.Getenv("SENDER_EMAIL")},
		Subject:  "Testing Emails",
		Host:     "smtp.gmail.com:587",
	}
	err = e2.SendEmail()
	if err != nil {
		t.Error(err)
	}
}
