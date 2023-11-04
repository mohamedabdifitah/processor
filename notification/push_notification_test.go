package notification

import (
	"fmt"
	"log"
	"os"
	"testing"

	"firebase.google.com/go/messaging"
)

func TestSendPushNotification(t *testing.T) {
	// This registration token comes from the client FCM SDKs.
	// testdevice := os.Getenv("TEST_DEVICE_ID")
	notification := messaging.Notification{
		Title: "Welcome",
		Body:  "Marwan welcome",
	}
	message := &messaging.Message{
		Data: map[string]string{
			"click_action": "FLUTTER_NOTIFICATION_CLICK",
			"sound":        "default",
			"status":       "done",
			"screen":       "screenB",
		},
		Notification: &notification,
		Token:        os.Getenv("TEST_DEVICE_ID"),
	}
	res, err := SendToastNotification(message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(res)
}
