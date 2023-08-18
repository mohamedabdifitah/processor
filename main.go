package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mohamedabdifitah/processor/controller"
	"github.com/mohamedabdifitah/processor/db"
	"github.com/mohamedabdifitah/processor/template"
)

func main() {
	if os.Getenv("APP_ENV") == "development" {
		err := godotenv.Load("local.env")
		if err != nil {
			log.Fatal(err)
		}
	}
	template.AllTemplates.LoadTemplates("./template/template.json")
	db.InitRedisClient()
	ListenTopic()
}
func ListenTopic() {
	topics := map[string]func([]byte){
		"new-order":                controller.NewOrderHandler,
		"order_accepted_resturant": controller.OrderAcceptedByResturantHandler,
	}
	topic := db.RedisClient.Subscribe(db.Ctx, "new-order", "order_accepted_resturant")
	channel := topic.Channel()
	for msg := range channel {
		topics[msg.Channel]([]byte(msg.Payload))
	}
}
