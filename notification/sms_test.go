package notification

import (
	"os"
	"testing"
)

func TestSendSms(t *testing.T) {
	sms := SMS{
		Message: "Hello world!",
		To:      os.Getenv("TEST_PHONE"),
	}
	sms.Send()
}
