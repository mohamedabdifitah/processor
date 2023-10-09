package notification

import (
	"os"
	"testing"

	"github.com/mohamedabdifitah/processor/utils"
)

func TestSendEmail(t *testing.T) {
	e := EmailMessage{
		Body:     "Helllo world\n",
		Receiver: []string{os.Getenv("SENDER_EMAIL")},
		Subject:  "Testing Emails",
		Host:     "smtp.gmail.com:587",
	}
	err := e.SendEmail()
	if err != nil {
		t.Error(err)
	}
	message, err := utils.AllTemplates.TempelateInjector(
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
