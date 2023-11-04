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
		"new_order":              service.HandleNewOrder,
		"merchant_accpted_order": service.HandleAcceptOrder,
		"order_canceled":         service.HandleCanceledOrder,
		"driver_accept_order":    service.HandleDriverAcceptOrder,
		"order_pickuped":         service.HandleOrderPickuped,
		"order_delivered":        service.HandleOrderDelivered,
		"order_dropped":          service.HandleDriverDropOrder,
		"order_ready":            service.HandleOrderReady,
		"order_preparing":        service.HandleOrderPreparing,
	}
	topic := pubsub.RedisClient.Subscribe(pubsub.Ctx,
		"merchant_accpted_order",
		"order_canceled",
		"order_preparing",
		"order_ready",
		"order_pickuped",
		"order_delivered",
		"new_order",
		"driver_accepted_order",
		"order_dropped",
	)
	channel := topic.Channel()
	for msg := range channel {
		topics[msg.Channel]([]byte(msg.Payload))
	}
}
