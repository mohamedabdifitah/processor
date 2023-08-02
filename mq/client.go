package broker

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var ch *amqp.Channel

func InitClient() {
	conn, err := amqp.Dial(os.Getenv("AMQP_URI"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err = conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	NewOrderQueue()
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
