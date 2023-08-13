package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mohamedabdifitah/processor/db"
	broker "github.com/mohamedabdifitah/processor/mq"
	"github.com/mohamedabdifitah/processor/template"
)

func main() {
	if os.Getenv("APP_ENV") == "development" {
		err := godotenv.Load("local.env")
		if err != nil {
			log.Fatal(err)
		}
	}
	template.AllTemplates.LoadTemplates("./template/text/text.json")
	go broker.InitProducer()
	db.InitRedisClient()
	broker.InitClient()
}
