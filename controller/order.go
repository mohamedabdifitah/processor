package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"firebase.google.com/go/messaging"
	"github.com/mohamedabdifitah/ecapi/db"
	"github.com/mohamedabdifitah/processor/notification"
	"github.com/mohamedabdifitah/processor/service"
	"github.com/mohamedabdifitah/processor/template"
)

type OrderBody struct {
	Drivers string   `json:"drivers"`
	Order   db.Order `json:"order"`
}

func NewOrderHandler(msg []byte) {
	var order db.Order
	// unmarshal message into order map
	if err := json.Unmarshal(msg, &order); err != nil {
		panic(err)
	}
	coordinates := order.PickUpLocation.Coordinates
	closestDrivers := service.SearchDrivers(5, coordinates[1], coordinates[0], 400, "m", true)
	if len(closestDrivers) == 0 {
		closestDrivers = service.SearchDrivers(5, coordinates[1], coordinates[0], 3000, "m", true)
	}
	// comma seperated ids e.g 12,2,3,4
	var stringids string
	for i, driver := range closestDrivers {
		id := strings.Split(driver.Name, ":")[1]
		if i == len(closestDrivers)-1 {
			stringids = stringids + id
			continue
		}
		stringids = stringids + id + ","

	}
	var body OrderBody = OrderBody{
		Drivers: stringids,
		Order:   order,
	}
	bodyjson, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}
	service.PublishTopic("socket-drives-order", bodyjson)
	resp, err := http.Get("http://localhost/driver/list?ids=" + stringids)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	writer := new(bytes.Buffer)
	byteResult, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	writer.Write(byteResult)
	var devices []*db.Device
	err = json.Unmarshal(byteResult, &devices)
	if err != nil {
		log.Fatal(err)
	}
	neworderTempalte, err := template.AllTemplates.TempelateInjector("NewOrder", map[string]string{
		"ResturantName": "Qoobey labadhagax",
		"from":          "Labadhagax",
		"to":            "Medina",
	})
	if err != nil {
		log.Fatal(err)
	}
	notificationbody := messaging.Notification{
		Title: fmt.Sprintf("New Order"),
		Body:  fmt.Sprintf(neworderTempalte),
	}
	// See documentation on defining a message payload.
	var times time.Duration = 6
	message := &messaging.Message{
		Data: map[string]string{
			"click_action": "FLUTTER_NOTIFICATION_CLICK",
			"sound":        "default",
			"status":       "done",
			"screen":       "orderScreen",
		},
		Notification: &notificationbody,
		Android: &messaging.AndroidConfig{
			TTL: &times,
		},
	}
	var list []string
	for _, v := range devices {
		list = append(list, v.DeviceId)
	}
	err = notification.SendMultipleNotifications(message, list)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("success message")
}
func OrderAcceptedByResturantHandler(msg []byte) {
}
