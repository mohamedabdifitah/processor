package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mohamedabdifitah/processor/pubsub"
	"github.com/mohamedabdifitah/processor/service"
	"github.com/mohamedabdifitah/processor/utils"
)

func main() {
	if os.Getenv("APP_ENV") == "development" {
		err := godotenv.Load(".env.local")
		if err != nil {
			log.Fatal(err)
		}
	}
	utils.AllTemplates.LoadTemplates("assets/json/template.json", "")
	pubsub.InitRedisClient()
	ListenTopic()
}
func ListenTopic() {
	topics := map[string]func([]byte){
		"new_order":                service.HandleNewOrder,
		"order_accepted_resturant": service.HandleAcceptOrder,
		"order_rejected_resturant": service.HandleMerchRejectOrder,
		"order_canceled":           service.HandleCanceledOrder,
		"driver_accept_order":      service.HandleDriverAcceptOrder,
		"order_pickuped":           service.HandleOrderPickuped,
		"order_delivered":          service.HandleOrderDelivered,
		"driver_drop_order":        service.HandleDriverDropOrder,
	}
	topic := pubsub.RedisClient.Subscribe(pubsub.Ctx, "new_order",
		"order_accepted_resturant",
		"order_rejected_resturant",
		"order_canceled",
		"driver_accept_order",
		"order_pickuped",
		"order_delivered",
		"driver_drop_order")
	channel := topic.Channel()
	for msg := range channel {
		topics[msg.Channel]([]byte(msg.Payload))
	}
}
