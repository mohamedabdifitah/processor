package service

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"firebase.google.com/go/messaging"
	notify "github.com/mohamedabdifitah/processor/notification"
	"github.com/mohamedabdifitah/processor/utils"
)

type Item struct {
	ItemExternalId string `json:"item_external_id"`
	Quantity       uint   `json:"quantity" `
	Price          uint   `json:"price" `
}
type Order struct {
	Id                    string    `json:"id,omitempty"`
	OrderValue            uint      `json:"order_value"  `
	Type                  string    `json:"type" `
	Items                 []Item    `json:"items"`
	DropOffPhone          string    `json:"dropoff_phone" `
	DropOffExteranlId     string    `json:"dropoff_external_id" `
	DropOffContactName    string    `json:"dropoff_contact_name" `
	DropOffTimeEstimated  int       `json:"dropoff_time_estimated" `
	DropOffAddress        string    `json:"dropoff_address" `
	DroOffLocation        Location  `json:"dropoff_location" `
	DropOffInstruction    string    `json:"dropoff_instructions" `
	Stage                 string    `json:"stage" `
	ActionIfUndeliverable string    `json:"action_if_undeliverable" `
	PickupAddress         string    `json:"pickup_address" `
	PickUpExternalId      string    `json:"pickup_external_id"`
	PickUpName            string    `json:"pickup_name"`
	PickUpPhone           string    `json:"pickup_phone"`
	PickUpLocation        Location  `json:"pickup_location"`
	PickupTime            time.Time `json:"pickup_time"`
	PickupEstimatedTime   int       `json:"pickup_estimated_time"`
	PickupReferenceTag    string    `json:"pickup_reference_tag" `
	DriverPhone           string    `json:"driver_phone"`
	DriverAllowedVehicles []string  `json:"driver_allowed_vehicles"  `
	DriverExternalId      string    `json:"driver_external_id"`
	Metadata              Metadata  `json:"metadata" `
	CancelReason          string    `json:"cancel_reason"`
}
type Metadata struct {
	CreatedAt time.Time `json:"created_at"`
}
type Location struct {
	Point       string    `json:"point"`
	Coordinates []float64 `json:"coordinates"`
}
type Device struct {
	DeviceId string ` json:"device_id"`
	Kind     string ` json:"kind"` // andriod ,ios
}

func HandleNewOrder(data []byte) {
	var order Order
	err := json.Unmarshal(data, &order)
	if err != nil {
		fmt.Println(err)
	}
	// merchant
	merchant := GetInformationMer(order.PickUpExternalId)
	notification := messaging.Notification{
		Title: fmt.Sprintf("Hey %s , New Order", merchant.BusinessName),
		Body:  fmt.Sprintf("New Order from %s please view items of order", order.DropOffContactName),
	}
	message := &messaging.Message{
		Data: map[string]string{
			"click_action": "FLUTTER_NOTIFICATION_CLICK",
			"sound":        "default",
			"status":       "done",
			"screen":       "",
		},
		Notification: &notification,
		Token:        merchant.Device.DeviceId,
	}
	_, err = notify.SendToastNotification(message)
	if err != nil {
		fmt.Println(err)
	}
	bytedata, err := json.Marshal(order)
	ok := notify.SendWebhook(merchant.Metadata.WebhookEndpoint+"/new/order", bytedata)
	if !ok {
		fmt.Println(err)
	}
	// drivers
	Drivers := GetDrivers(order)

	for _, driver := range Drivers {
		pickupDist := utils.CalculateDistance(order.PickUpLocation.Coordinates[0], order.PickUpLocation.Coordinates[1], driver.Location.Coordinates[0], driver.Location.Coordinates[1], "")
		pickupDistance := utils.Converter(pickupDist)
		dropoffdist := utils.CalculateDistance(order.DroOffLocation.Coordinates[0], order.DroOffLocation.Coordinates[1], driver.Location.Coordinates[0], driver.Location.Coordinates[1], "")
		dropoffistance := utils.Converter(dropoffdist)
		template, err := utils.AllTemplates.TempelateInjector("NewOrder", map[string]string{
			"merchantname":    order.PickUpName,
			"numberOfItems":   fmt.Sprintf("%s", strconv.Itoa(len(order.Items))),
			"pickupdistance":  pickupDistance,
			"pickupAddress":   order.PickupAddress,
			"dropoffdistance": dropoffistance,
			"dropoffAddress":  order.DropOffAddress,
		})
		if err != nil {
			fmt.Println(err)
		}
		notification = messaging.Notification{
			Title: fmt.Sprintf("New Delivery"),
			Body:  template,
		}
		message = &messaging.Message{
			Data: map[string]string{
				"click_action": "FLUTTER_NOTIFICATION_CLICK",
				"sound":        "default",
				"status":       "done",
				"screen":       "",
			},
			Notification: &notification,
			Token:        driver.Device.DeviceId,
		}
		_, err = notify.SendToastNotification(message)
		if err != nil {
			fmt.Println(err)
		}
		bytedata, err := json.Marshal(order)
		ok := notify.SendWebhook(driver.Metadata.WebhookEndpoint+"/new/order", bytedata)
		if !ok {
			fmt.Println(err)
		}
	}
}
func HandleAcceptOrder(data []byte) {
	var order Order
	err := json.Unmarshal(data, &order)
	if err != nil {
		fmt.Println(err)
	}
	// customer
	customer := GetInformationCustomer(order.DropOffExteranlId)
	notification := messaging.Notification{
		Title: "Order is Accepted",
		Body:  fmt.Sprintf("your order for %s is accepted", order.PickUpName),
	}
	message := &messaging.Message{
		Data: map[string]string{
			"click_action": "FLUTTER_NOTIFICATION_CLICK",
			"sound":        "default",
			"status":       "done",
			"screen":       "",
		},
		Notification: &notification,
		Token:        customer.Device.DeviceId,
	}
	_, err = notify.SendToastNotification(message)
	if err != nil {
		fmt.Println(err)
	}
	bytedata, err := json.Marshal(order)
	ok := notify.SendWebhook(customer.Metadata.WebhookEndpoint+"/merchant/accept", bytedata)
	if !ok {
		fmt.Println(err)
	}
}
func HandleDriverAcceptOrder(data []byte) {
	var order Order
	err := json.Unmarshal(data, &order)
	if err != nil {
		fmt.Println(err)
	}
	driver := GetInformationDriver(order.DriverExternalId)
	// customer
	customer := GetInformationCustomer(order.DropOffExteranlId)
	notification := messaging.Notification{
		Title: "Driver accept your order",
		Body:  fmt.Sprintf("your order for %s is accepted by driver %s", order.PickUpName, driver.GivenName),
	}
	message := &messaging.Message{
		Data: map[string]string{
			"click_action": "FLUTTER_NOTIFICATION_CLICK",
			"sound":        "default",
			"status":       "done",
			"screen":       "",
		},
		Notification: &notification,
		Token:        customer.Device.DeviceId,
	}
	_, err = notify.SendToastNotification(message)
	if err != nil {
		fmt.Println(err)
	}
	bytedata, err := json.Marshal(order)
	ok := notify.SendWebhook(customer.Metadata.WebhookEndpoint+"/driver/accept", bytedata)
	if !ok {
		fmt.Println(err)
	}
}
func HandleCanceledOrder(data []byte) {
	var order Order
	err := json.Unmarshal(data, &order)
	if err != nil {
		fmt.Println(err)
	}

	reason := strings.Split(order.CancelReason, " ")[0]
	if reason == "CANCEL_FROM_MERCHANT" {
		customer := GetInformationCustomer(order.DropOffExteranlId)
		notification := messaging.Notification{
			Title: "Order canceled",
			Body:  fmt.Sprintf("your order from %s is canceled due to %s", order.PickUpName, strings.Split(order.CancelReason, " ")[1]),
		}
		message := &messaging.Message{
			Data: map[string]string{
				"click_action": "FLUTTER_NOTIFICATION_CLICK",
				"sound":        "default",
				"status":       "done",
				"screen":       "",
			},
			Notification: &notification,
			Token:        customer.Device.DeviceId,
		}
		_, err = notify.SendToastNotification(message)
		if err != nil {
			fmt.Println(err)
		}
		bytedata, err := json.Marshal(order)
		ok := notify.SendWebhook(customer.Metadata.WebhookEndpoint+"/order/cancel", bytedata)
		if !ok {
			fmt.Println(err)
		}
	} else if reason == "CANCEL_FROM_CUSTOMER" {
		// merchant
		merchant := GetInformationMer(order.PickUpExternalId)
		notification := messaging.Notification{
			Title: "Order canceled",
			Body:  fmt.Sprint("%s cancels his order", order.DropOffContactName),
		}
		message := &messaging.Message{
			Data: map[string]string{
				"click_action": "FLUTTER_NOTIFICATION_CLICK",
				"sound":        "default",
				"status":       "done",
				"screen":       "",
			},
			Notification: &notification,
			Token:        merchant.Device.DeviceId,
		}
		_, err = notify.SendToastNotification(message)
		if err != nil {
			fmt.Println(err)
		}
		bytedata, err := json.Marshal(order)
		ok := notify.SendWebhook(merchant.Metadata.WebhookEndpoint+"/order/cancel", bytedata)
		if !ok {
			fmt.Println(err)
		}
	}
}
func HandleDriverDropOrder(data []byte) {
	var order Order
	err := json.Unmarshal(data, &order)
	if err != nil {
		fmt.Println(err)
	}
	driver := GetInformationDriver(order.DriverExternalId)
	// merchant
	merchant := GetInformationMer(order.PickUpExternalId)
	notification := messaging.Notification{
		Title: "Driver cancels delivery",
		Body:  fmt.Sprint("%s cancels his order of %s", driver.GivenName, order.DropOffContactName),
	}
	message := &messaging.Message{
		Data: map[string]string{
			"click_action": "FLUTTER_NOTIFICATION_CLICK",
			"sound":        "default",
			"status":       "done",
			"screen":       "",
		},
		Notification: &notification,
		Token:        merchant.Device.DeviceId,
	}
	_, err = notify.SendToastNotification(message)
	if err != nil {
		fmt.Println(err)
	}
	bytedata, err := json.Marshal(order)
	ok := notify.SendWebhook(merchant.Metadata.WebhookEndpoint+"/order/cancel", bytedata)
	if !ok {
		fmt.Println(err)
	}
	// drivers
	Drivers := GetDrivers(order)

	for _, driver := range Drivers {
		pickupDist := utils.CalculateDistance(order.PickUpLocation.Coordinates[0], order.PickUpLocation.Coordinates[1], driver.Location.Coordinates[0], driver.Location.Coordinates[1], "")
		pickupDistance := utils.Converter(pickupDist)
		dropoffdist := utils.CalculateDistance(order.DroOffLocation.Coordinates[0], order.DroOffLocation.Coordinates[1], driver.Location.Coordinates[0], driver.Location.Coordinates[1], "")
		dropoffistance := utils.Converter(dropoffdist)
		template, err := utils.AllTemplates.TempelateInjector("NewOrder", map[string]string{
			"merchantname":    order.PickUpName,
			"numberOfItems":   fmt.Sprintf("%s", strconv.Itoa(len(order.Items))),
			"pickupdistance":  pickupDistance,
			"pickupAddress":   order.PickupAddress,
			"dropoffdistance": dropoffistance,
			"dropoffAddress":  order.DropOffAddress,
		})
		if err != nil {
			fmt.Println(err)
		}
		notification := messaging.Notification{
			Title: fmt.Sprintf("New Delivery"),
			Body:  template,
		}
		message := &messaging.Message{
			Data: map[string]string{
				"click_action": "FLUTTER_NOTIFICATION_CLICK",
				"sound":        "default",
				"status":       "done",
				"screen":       "",
			},
			Notification: &notification,
			Token:        driver.Device.DeviceId,
		}
		_, err = notify.SendToastNotification(message)
		if err != nil {
			fmt.Println(err)
		}
		bytedata, err := json.Marshal(order)
		ok := notify.SendWebhook(driver.Metadata.WebhookEndpoint+"/new/order", bytedata)
		if !ok {
			fmt.Println(err)
		}
	}
}
func HandleOrderDelivered(data []byte) {
	var order Order
	err := json.Unmarshal(data, &order)
	if err != nil {
		fmt.Println(err)
	}
	merchant := GetInformationMer(order.PickUpExternalId)
	notification := messaging.Notification{
		Title: "Order is delivered",
		Body:  fmt.Sprint("%s order is  delivered", order.DropOffContactName),
	}
	message := &messaging.Message{
		Data: map[string]string{
			"click_action": "FLUTTER_NOTIFICATION_CLICK",
			"sound":        "default",
			"status":       "done",
			"screen":       "",
		},
		Notification: &notification,
		Token:        merchant.Device.DeviceId,
	}
	_, err = notify.SendToastNotification(message)
	if err != nil {
		fmt.Println(err)
	}
	bytedata, err := json.Marshal(order)
	ok := notify.SendWebhook(merchant.Metadata.WebhookEndpoint+"/order/delivered", bytedata)
	if !ok {
		fmt.Println(err)
	}
	// customer
	driver := GetInformationDriver(order.DriverExternalId)
	customer := GetInformationCustomer(order.DropOffExteranlId)
	notification = messaging.Notification{
		Title: "Order is delivered",
		Body:  fmt.Sprintf("your driver %s dropped your order, please check your order. if you have trouble finding something please call us %s ", driver.GivenName, os.Getenv("TEST_PHONE")),
	}
	message = &messaging.Message{
		Data: map[string]string{
			"click_action": "FLUTTER_NOTIFICATION_CLICK",
			"sound":        "default",
			"status":       "done",
			"screen":       "",
		},
		Notification: &notification,
		Token:        customer.Device.DeviceId,
	}
	_, err = notify.SendToastNotification(message)
	if err != nil {
		fmt.Println(err)
	}
	bytedata, err = json.Marshal(order)
	ok = notify.SendWebhook(customer.Metadata.WebhookEndpoint+"/order/delivered", bytedata)
	if !ok {
		fmt.Println(err)
	}
}
func HandleOrderPickuped(data []byte) {
	var order Order
	err := json.Unmarshal(data, &order)
	if err != nil {
		fmt.Println(err)
	}
	// customer
	driver := GetInformationDriver(order.DriverExternalId)
	customer := GetInformationCustomer(order.DropOffExteranlId)
	notification := messaging.Notification{
		Title: "Order is pickuped",
		Body:  fmt.Sprintf("your driver %s pick up your order it will reach you %d minutes", driver.GivenName, order.DropOffTimeEstimated/60),
	}
	message := &messaging.Message{
		Data: map[string]string{
			"click_action": "FLUTTER_NOTIFICATION_CLICK",
			"sound":        "default",
			"status":       "done",
			"screen":       "",
		},
		Notification: &notification,
		Token:        customer.Device.DeviceId,
	}
	_, err = notify.SendToastNotification(message)
	if err != nil {
		fmt.Println(err)
	}
	bytedata, err := json.Marshal(order)
	ok := notify.SendWebhook(customer.Metadata.WebhookEndpoint+"/order/delivered", bytedata)
	if !ok {
		fmt.Println(err)
	}
}
func HandleOrderReady(data []byte) {
	var order Order
	err := json.Unmarshal(data, &order)
	if err != nil {
		fmt.Println(err)
	}
	// customer
	driver := GetInformationDriver(order.DriverExternalId)
	notification := messaging.Notification{
		Title: "Order is ready to delivered",
		Body:  fmt.Sprintf("%s's Order is ready for pickup at %s. Go ahead and grab your order - make sure you check items as usual", order.DropOffContactName, order.PickUpName),
	}
	message := &messaging.Message{
		Data: map[string]string{
			"click_action": "FLUTTER_NOTIFICATION_CLICK",
			"sound":        "default",
			"status":       "done",
			"screen":       "",
		},
		Notification: &notification,
		Token:        driver.Device.DeviceId,
	}
	_, err = notify.SendToastNotification(message)
	if err != nil {
		fmt.Println(err)
	}
	sms := notify.SMS{
		To:      order.DriverPhone,
		Message: fmt.Sprintf("%s's Order is ready for pickup at %s. Go ahead and grab your order - make sure you check items as usual", order.DropOffContactName, order.PickUpName),
	}
	err = sms.Send()
	if err != nil {
		fmt.Println(err)
	}
	bytedata, err := json.Marshal(order)
	ok := notify.SendWebhook(driver.Metadata.WebhookEndpoint+"/order/ready", bytedata)
	if !ok {
		fmt.Println(err)
	}
}
func HandleOrderPreparing(data []byte) {
	var order Order
	err := json.Unmarshal(data, &order)
	if err != nil {
		fmt.Println(err)
	}

	// customer
	driver := GetInformationDriver(order.DriverExternalId)
	if err != nil {
		fmt.Println(err)
	}
	notification := messaging.Notification{
		Title: "Order is processing",
		Body:  fmt.Sprintf("%s's Order is almost ready for pickup at %d minutes", order.DropOffContactName, order.PickupEstimatedTime/60),
	}
	message := &messaging.Message{
		Data: map[string]string{
			"click_action": "FLUTTER_NOTIFICATION_CLICK",
			"sound":        "default",
			"status":       "done",
			"screen":       "",
		},
		Notification: &notification,
		Token:        driver.Device.DeviceId,
	}
	_, err = notify.SendToastNotification(message)
	if err != nil {
		fmt.Println(err)
	}
	sms := notify.SMS{
		To:      order.DriverPhone,
		Message: fmt.Sprintf("%s's Order is almost ready for pickup at %d minutes", order.DropOffContactName, order.PickupEstimatedTime/60),
	}
	err = sms.Send()
	if err != nil {
		fmt.Println(err)
	}
	bytedata, err := json.Marshal(order)
	ok := notify.SendWebhook(driver.Metadata.WebhookEndpoint+"/order/preparing", bytedata)
	if !ok {
		fmt.Println(err)
	}
}
