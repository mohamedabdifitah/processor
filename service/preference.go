package service

import (
	"os"

	"firebase.google.com/go/messaging"
	"github.com/mohamedabdifitah/processor/notification"
	notify "github.com/mohamedabdifitah/processor/notification"
)

func SendPushNotificationPlatforms(token string, page string, title string, body string) error {
	notification := messaging.Notification{
		Title: title,
		Body:  body,
	}
	message := &messaging.Message{
		Data: map[string]string{
			"click_action": "FLUTTER_NOTIFICATION_CLICK",
			"sound":        "default",
			"status":       "done",
			"screen":       page,
		},
		Notification: &notification,
		Token:        token,
	}
	_, err := notify.SendToastNotification(message)
	if err != nil {
		return err
	}
	return nil
}
func SendWebsocketNotification(event string, body interface{}, rec string, socket notify.Socket) error {
	var message notify.Event = notify.Event{
		Name:    event,
		Message: body,
	}
	err := socket.Send(socket.Clients[rec], message)
	if err != nil {
		return err
	}
	return nil
}
func SendNotifictionSms(message string, to string) error {
	sms := notification.SMS{
		Message: message,
		To:      to,
	}
	err := sms.Send()
	if err != nil {
		return err
	}
	return nil
}

func SendEmailNotification(message string, subject string, rec ...string) error {
	e := notification.EmailMessage{
		Body:     message,
		Mime:     "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\"",
		Receiver: rec,
		Subject:  subject,
		Host:     os.Getenv("SENDER_EMAIL_HOST_PORT"),
	}
	err := e.SendEmail()
	if err != nil {
		return err
	}
	return nil
}
