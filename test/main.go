package main

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	// service.failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	// fmt.Printf("Failed to publish a message" + err.Error())
	fmt.Println("connect to RabbitMQ")
	// Create a channel
	ch, err := conn.Channel()
	// service.failOnError(err, "Failed to open a channel")
	defer ch.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = ch.PublishWithContext(
		ctx,
		"",               // exchange
		"reset_password", // routing key (queue name)
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte{},
		})
	if err != nil {
		fmt.Printf("Failed to publish a message" + err.Error())
	}
}
