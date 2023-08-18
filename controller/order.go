package controller

import (
	"encoding/json"
	"fmt"

	"github.com/mohamedabdifitah/processor/service"
)

type Order struct {
	PickupLocation   Location `json:"pickup_location"`
	PickupExternalId string   `json:"pickup_external_id"`
	PickUpPhone      string   `json:"pickup_phone"`
}
type Location struct {
	Coordinates []float64 `json:"coordinates"`
	Point       string    `json:"point"`
}

func NewOrderHandler(msg []byte) {
	var order Order
	// unmarshal message into order map
	if err := json.Unmarshal(msg, &order); err != nil {
		panic(err)
	}
	coordinates := order.PickupLocation.Coordinates
	drivers := service.SearchDrivers(20, coordinates[1], coordinates[0], 400, "m", true)
	if len(drivers) == 0 {
		drivers = service.SearchDrivers(20, coordinates[1], coordinates[0], 20000, "m", true)
	}
	fmt.Println(drivers)
}
func OrderAcceptedByResturantHandler(msg []byte) {
}
