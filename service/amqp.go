package service

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s", msg, err)
		return
	}
}

var Channel *amqp.Channel
var Queues map[string]func(amqp.Delivery) = map[string]func(amqp.Delivery){
	"new_order":              HandleNewOrder,
	"order_ready":            HandleOrderReady,
	"order_dropped":          HandleDriverDropOrder,
	"order_canceled":         HandleCanceledOrder,
	"order_pickuped":         HandleOrderPickuped,
	"reset_password":         HandleResetPassword,
	"order_preparing":        HandleOrderPreparing,
	"order_delivered":        HandleOrderDelivered,
	"driver_accepted_order":  HandleDriverAcceptOrder,
	"order_ready_assigment":  HandleOrderAssignment,
	"merchant_accpted_order": HandleAcceptOrder,
	"verification":           HandleVerification,
}

func InitAmqp() {
	// Connect to RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	fmt.Println("connect to RabbitMQ")
	// Create a channel
	Channel, err = conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer Channel.Close()

	// Declare multiple queues

	for queueName, handler := range Queues {
		q, err := Channel.QueueDeclare(
			queueName, // name
			true,      // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		failOnError(err, "Failed to declare a queue")

		// Start a consumer for each queue
		go consumeQueue(q.Name, handler)
	}
	// Run indefinitely to keep the program running
	select {}
}

func consumeQueue(queueName string, handler func(amqp.Delivery)) {
	messages, err := Channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	failOnError(err, "Failed to register a consumer")

	for message := range messages {
		handler(message)
	}
}
func ProduceMessages(queueName string, body []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := Channel.PublishWithContext(
		ctx,
		"",        // exchange
		queueName, // routing key (queue name)
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		fmt.Printf("Failed to publish a message" + err.Error())
	}
}
