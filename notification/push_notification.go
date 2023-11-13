package notification

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

func SendMultipleNotifications(message *messaging.Message, rec []string) error {
	if len(rec) == 0 {
		return fmt.Errorf("there are no recievers")
	}
	for _, v := range rec {
		message.Token = v
		_, err := SendToastNotification(message)
		if err != nil {
			return err
		}
	}
	return nil
}

// send muluticast
func SendToastNotification(message *messaging.Message) (string, error) {
	wd, _ := os.Getwd()
	path := filepath.Join(wd + "/assets/json/serviceAccountKey.json")
	if filepath.Base(wd) == "notification" {
		path = filepath.Join(filepath.Dir(wd) + "/assets/json/serviceAccountKey.json")
	}
	opt := option.WithCredentialsFile(path)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return "", fmt.Errorf("error initializing app: %v \n ", err)
	}
	// Obtain a messaging.Client from the App.
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting messaging client: %v", err)
	}
	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := client.Send(ctx, message)
	if err != nil {
		return "", err
	}
	return response, nil
}
