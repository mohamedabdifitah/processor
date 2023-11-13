package service

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	notify "github.com/mohamedabdifitah/processor/notification"
	amqp "github.com/rabbitmq/amqp091-go"
)

func HandleNewOrder(data amqp.Delivery) {
	var order Order
	err := json.Unmarshal(data.Body, &order)
	if err != nil {
		fmt.Println(err)
	}
	// merchant
	merchant := GetInformationMer(order.PickUpExternalId)
	var title string = "New Order"
	var body string = fmt.Sprintf("New Order from %s please view items of order", order.DropOffContactName)
	var link string = "order/" + order.Id
	// push notification
	SendPushNotificationPlatforms(merchant.Device.DeviceId, link, title, body)
	// websocket notification
	SendWebsocketNotification("new_order", order, merchant.Id, notify.MerchantSocket)
	// webhook notification
	bytedata, err := json.Marshal(order)
	ok := notify.SendWebhook(merchant.Metadata.WebhookEndpoint+"/new/order", bytedata)
	if !ok {
		fmt.Println(err)
	}
	ProduceMessages("order_ready_assigment", bytedata)
}
func HandleAcceptOrder(data amqp.Delivery) {
	var order Order
	err := json.Unmarshal(data.Body, &order)
	if err != nil {
		fmt.Println(err)
	}
	customer := GetInformationCustomer(order.DropOffExteranlId)
	var title string = order.PickUpName + " Accepted yout order."
	var body string = "please review or follow order activity"
	var link string = "/activity/order/" + order.Id
	// push notification
	SendPushNotificationPlatforms(customer.Device.DeviceId, link, title, body)
	// websocket notification
	SendWebsocketNotification("order_status", order, customer.Id, notify.ClientSocket)
}
func HandleDriverAcceptOrder(data amqp.Delivery) {
	var order Order
	err := json.Unmarshal(data.Body, &order)
	if err != nil {
		fmt.Println(err)
	}
	driver := GetInformationDriver(order.DriverExternalId)
	// customer
	customer := GetInformationCustomer(order.DropOffExteranlId)
	var title string = "Driver accept your order"
	var body string = fmt.Sprintf("your order for %s is accepted by driver %s , please review or follow order activity ", order.PickUpName, driver.GivenName)
	var link string = "/activity/order/" + order.Id
	// push notification
	SendPushNotificationPlatforms(customer.Device.DeviceId, link, title, body)
	// websocket notification
	SendWebsocketNotification("order_status", order, customer.Id, notify.ClientSocket)
	merchant := GetInformationMer(order.PickUpExternalId)
	title = "Driver accepted order " + order.DisplayId
	body = fmt.Sprintf("your order from %s is accepted by driver %s , please review or follow order activity ", order.DropOffContactName, driver.GivenName)
	link = "/activity/order/" + order.Id
	// push notification
	SendPushNotificationPlatforms(merchant.Device.DeviceId, link, title, body)
	// websocket notification
	SendWebsocketNotification("order_status", order, merchant.Id, notify.MerchantSocket)
}
func HandleCanceledOrder(data amqp.Delivery) {
	var order Order
	err := json.Unmarshal(data.Body, &order)
	if err != nil {
		fmt.Println(err)
	}

	reason := strings.Split(order.CancelReason, " ")[0]
	var driver Driver
	var hasDriver = false
	if order.DriverExternalId != "" {
		driver = GetInformationDriver(order.DriverExternalId)
		hasDriver = true
	}
	if reason == "CANCEL_FROM_MERCHANT" {

		customer := GetInformationCustomer(order.DropOffExteranlId)

		var message string = "your order to " + order.PickUpName + " is canceled by merchant due to " + strings.Split(order.CancelReason, " ")[1] + " please open the app and order from new merchant"
		SendNotifictionSms(message, customer.Phone)
		var title string = "Order is cancelled"
		var body string = "your order to " + order.PickUpName + " is cancelled due to " + strings.Split(order.CancelReason, " ")[1]
		var link string = "/activity/order/" + order.Id
		SendPushNotificationPlatforms(customer.Device.DeviceId, link, title, body)
		if hasDriver {
			SendWebsocketNotification("order_canceled", order, driver.Id, notify.DriverSocket)
			var title string = "Order is cancelled"
			var body string = "your order from " + order.PickUpName + " to " + order.DropOffContactName + "is cancelled due to " + strings.Split(order.CancelReason, " ")[1]
			var link string = "/delivering/order/" + order.Id
			SendPushNotificationPlatforms(driver.Device.DeviceId, link, title, body)
		}

	} else if reason == "CANCEL_FROM_CUSTOMER" {
		// merchant
		merchant := GetInformationMer(order.PickUpExternalId)
		var title string = "Order" + order.DisplayId + "is canceled"
		var body string = order.DropOffContactName + "cancels his order , please do not fulfil this order " + order.DisplayId + " or click this notification"
		var link string = "order/" + order.Id
		SendPushNotificationPlatforms(merchant.Device.DeviceId, link, title, body)
		SendWebsocketNotification("order_canceled", order, merchant.Id, notify.MerchantSocket)
		bytedata, err := json.Marshal(order)
		if err != nil {
			fmt.Println(err)
		}
		ok := notify.SendWebhook(merchant.Metadata.WebhookEndpoint+"/order/cancel", bytedata)
		if !ok {
			fmt.Println(err)
		}
		if hasDriver {
			SendWebsocketNotification("order_canceled", order, driver.Id, notify.DriverSocket)
			var title string = "Order is cancelled"
			var body string = "your order from " + order.PickUpName + " to " + order.DropOffContactName + "is cancelled due to " + strings.Split(order.CancelReason, " ")[1]
			var link string = "/delivering/order/" + order.Id
			SendPushNotificationPlatforms(driver.Device.DeviceId, link, title, body)
		}
	}
}
func HandleDriverDropOrder(data amqp.Delivery) {
	var order Order
	err := json.Unmarshal(data.Body, &order)
	if err != nil {
		fmt.Println(err)
	}
	driver := GetInformationDriver(order.DriverExternalId)
	// merchant
	merchant := GetInformationMer(order.PickUpExternalId)
	var title string = "Driver cancels delivery"
	var body string = driver.GivenName + "cancels his order of" + order.DisplayId + "to " + order.DropOffContactName
	var link string = "order/" + order.Id
	SendPushNotificationPlatforms(merchant.Device.DeviceId, link, title, body)
	SendWebsocketNotification("driver_dropp_order", order, merchant.Id, notify.MerchantSocket)
	bytedata, err := json.Marshal(order)
	ok := notify.SendWebhook(merchant.Metadata.WebhookEndpoint+"/order/cancel", bytedata)
	if !ok {
		fmt.Println(err)
	}
	ProduceMessages("order_ready_assigment", bytedata)
}
func HandleOrderDelivered(data amqp.Delivery) {
	var order Order
	err := json.Unmarshal(data.Body, &order)
	if err != nil {
		fmt.Println(err)
	}
	merchant := GetInformationMer(order.PickUpExternalId)
	driver := GetInformationDriver(order.DriverExternalId)
	var title string = "Order " + order.DisplayId + " is delivered"
	var body string = driver.GivenName + " delivered his order to " + order.DropOffContactName
	var link string = "order/" + order.Id
	SendPushNotificationPlatforms(merchant.Device.DeviceId, link, title, body)
	bytedata, err := json.Marshal(order)
	ok := notify.SendWebhook(merchant.Metadata.WebhookEndpoint+"/order/delivered", bytedata)
	if !ok {
		fmt.Println(err)
	}
	// customer
	customer := GetInformationCustomer(order.DropOffExteranlId)
	title = "Order is completed "
	body = fmt.Sprintf("your driver %s dropped your order, please check your order. if you have trouble finding something please call us %s ", driver.GivenName, os.Getenv("TEST_PHONE"))
	link = "activity/order/" + order.Id
	SendPushNotificationPlatforms(customer.Device.DeviceId, link, title, body)
	SendNotifictionSms(title+body, customer.Phone)
}
func HandleOrderPickuped(data amqp.Delivery) {
	var order Order
	err := json.Unmarshal(data.Body, &order)
	if err != nil {
		fmt.Println(err)
	}
	// customer
	driver := GetInformationDriver(order.DriverExternalId)
	customer := GetInformationCustomer(order.DropOffExteranlId)
	var title string = "Order is pickuped"
	var body string = fmt.Sprintf("your driver %s pick up your order it will reach you %d minutes", driver.GivenName, order.DropOffTimeEstimated/60)
	var link string = "order/id/" + order.Id
	SendPushNotificationPlatforms(customer.Device.DeviceId, link, title, body)
	var message string = fmt.Sprintf("your driver %s pick up your order it will reach you %d minutes", driver.GivenName, order.DropOffTimeEstimated/60)
	SendNotifictionSms(message, customer.Phone)
}
func HandleOrderReady(data amqp.Delivery) {
	var order Order
	err := json.Unmarshal(data.Body, &order)
	if err != nil {
		fmt.Println(err)
	}
	// customer
	driver := GetInformationDriver(order.DriverExternalId)
	var title string = "Order is ready to Pick up"
	var body string = fmt.Sprintf("%s's Order is ready for pickup at %s. Go ahead and grab your order - make sure you check items as usual", order.DropOffContactName, order.PickUpName)
	var link string = "/activity/order/" + order.Id
	SendPushNotificationPlatforms(driver.Device.DeviceId, link, title, body)
	var message string = fmt.Sprintf("%s's Order is ready for pickup at %s. Go ahead and grab your order - make sure you check items as usual", order.DropOffContactName, order.PickUpName)
	SendNotifictionSms(message, order.DriverPhone)
}
func HandleOrderPreparing(data amqp.Delivery) {
	var order Order
	err := json.Unmarshal(data.Body, &order)
	if err != nil {
		fmt.Println(err)
	}

	// driver
	if err != nil {
		fmt.Println(err)
	}
	var message string = order.PickUpName + "start preparing order , is almost ready for pickup at " + fmt.Sprintf("%d", order.PickupEstimatedTime/60) + "minutes"
	SendNotifictionSms(message, order.DriverPhone)
}
