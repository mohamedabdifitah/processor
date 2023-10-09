package notification

import (
	"context"
	"fmt"

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
func SendToastNotification(message *messaging.Message) (*string, error) {
	opt := option.WithCredentialsFile("unknown")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v\n", err)
	}
	// Obtain a messaging.Client from the App.
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting messaging client: %v", err)
	}
	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := client.Send(ctx, message)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
