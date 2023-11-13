package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mohamedabdifitah/processor/notification"
	"github.com/mohamedabdifitah/processor/service"
)

func main() {
	if os.Getenv("APP_ENV") == "development" {
		err := godotenv.Load(".env.local")
		if err != nil {
			log.Fatal(err)
		}
	}
	service.InitRedisClient()
	go notification.Initsocket()
	service.InitAmqp()
}
