package notification

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"firebase.google.com/go/messaging"
	"github.com/joho/godotenv"
)

func TestSendPushNotification(t *testing.T) {
	err := godotenv.Load("../local.env")
	if err != nil {
		log.Fatal(err)
	}
	// This registration token comes from the client FCM SDKs.
	registrationToken := os.Getenv("TEST_DEVICE_ID")
	notification := messaging.Notification{
		Title: "Welcome",
		Body:  "welcome to the messaging ",
	}
	// See documentation on defining a message payload.
	var times time.Duration = 6
	message := &messaging.Message{
		Data: map[string]string{
			"click_action": "FLUTTER_NOTIFICATION_CLICK",
			"sound":        "default",
			"status":       "done",
			"screen":       "screenB",
		},
		Notification: &notification,
		Token:        registrationToken,
		Android: &messaging.AndroidConfig{
			TTL: &times,
		},
	}
	res, err := SendToastNotification(message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(res)

}
