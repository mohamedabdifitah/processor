package controller

import (
	"github.com/rabbitmq/amqp091-go"
)

func NewOrderHandler(msg amqp091.Delivery) {
}
func OrderAcceptedByResturantHandler(msg amqp091.Delivery) {
}
func OrderRejectedByResturantHandler(msg amqp091.Delivery) {
}
func OrderAccpetedByDriverHandler(msg amqp091.Delivery) {
}
func OrderReadyHandler(msg amqp091.Delivery) {
}
func OrderPickupedHandlers(msg amqp091.Delivery) {
}
