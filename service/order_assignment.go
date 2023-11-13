package service

import (
	"encoding/json"
	"fmt"
	"time"

	// "github.com/mohamedabdifitah/processor/internal"
	notify "github.com/mohamedabdifitah/processor/notification"
	amqp "github.com/rabbitmq/amqp091-go"
)

func AsignOrderToDriver(order Order, driver Driver) {
	// push notification
	var title string = "New Order GO To" + order.PickUpName
	var body string = fmt.Sprintf("New Order from %s please view items of order", order.DropOffContactName)
	var invitationlink string = "invitation/" + order.Id
	// push notification
	err := SendPushNotificationPlatforms(driver.Device.DeviceId, invitationlink, title, body)
	if err != nil {
		fmt.Println(err)
	}
	// websocket notification
	err = SendWebsocketNotification("invitation", order, driver.Id, notify.DriverSocket)
	if err != nil {
		fmt.Println(err)
	}
	// invitation time to live
	RedisClient.Set(Ctx, "inviation:"+order.Id, driver.Id, 60)

}
func checkOrderState(key string) bool {
	state := RedisClient.Get(Ctx, key).Val()
	return state == "accepted"
}
func HandleOrderAssignment(msg amqp.Delivery) {
	var order Order
	err := json.Unmarshal(msg.Body, &order)
	if err != nil {
		fmt.Println(err)
	}
	if checkOrderState("order:" + order.Id) {
		return
	}
	var key int = 0
	drivers := FindBestDriver(order)
	AsignOrderToDriver(order, drivers[key])
	var wholetime int = 10
	if len(drivers) < 8 {
		wholetime = len(drivers) * 60
	}
	timeout := time.After(time.Duration(wholetime) * time.Second)
	ticker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-timeout:
			ProduceMessages("order_invitation_timeout", msg.Body)
			return
		case <-ticker.C:
			if checkOrderState("order:" + order.Id) {
				return
			}
			AsignOrderToDriver(order, drivers[key])
			key = key + 1
		}
	}
}
