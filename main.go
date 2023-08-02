package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	broker "github.com/mohamedabdifitah/processor/mq"
)

func main() {
	if os.Getenv("APP_ENV") == "development" {
		err := godotenv.Load("local.env")
		if err != nil {
			log.Fatal(err)
		}
	}
	go broker.InitProducer()

	broker.InitClient()
}
