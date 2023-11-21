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
		Title:    "Test Notifications",
		Body:     "Firebase  Order test 14",
		ImageURL: "https://s3-alpha-sig.figma.com/img/939f/d722/70ffd81364346efa0b542910aed37b59?Expires=1701043200&Signature=mAKkqAesjP8BqYVJZAfD4MezlZG3LybOvB9zmin4uNwNT8G-ksMtoEJvctUHXt7XB~7YYJnkaptbbM4ul2XPqQxzKfoPQkWYMvx5UegS5RTwk8tiUl4SAR8A1IMDbxfZIFEpkYtg1iTtXX97MLwZ5ZMVaYwfvR8R~XxRZonfmbIOG3HZxXpm2KBepnoFdT3FoXCpvv-fVmZMcWLjpYCtVaGiW2lvsan8Qw69TXp236A6FQIhvy7yGT8dYdEvlXUVc4nfVqAYfP~mMCCWMf9UZBIfkmVw79TtGtXPr0tl6L8QxmibjPQhozwDcYM6ThAz2u1okU9NzSzDc6giz9gQYQ__&Key-Pair-Id=APKAQ4GOSFWCVNEHN3O4",
	}
	message := &messaging.Message{
		Data: map[string]string{
			// "click_action": "FLUTTER_NOTIFICATION_CLICK",
			"sound":  "default",
			"status": "done",
			"screen": "/dashboard/search",
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
