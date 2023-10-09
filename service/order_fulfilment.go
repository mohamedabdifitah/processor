package service

import (
	"encoding/json"
	"log"
	"strings"
	"time"
)

type Order struct {
	Id                    string    `json:"id,omitempty"`
	OrderValue            uint      `json:"order_value"  `
	Type                  string    `json:"type" `
	DropOffPhone          string    `json:"dropoff_phone" `
	DropOffExteranlId     string    `json:"dropoff_external_id" `
	DropOffContactName    string    `json:"dropoff_contact_name" `
	DropOffTimeEstimated  time.Time `json:"dropoff_time_estimated" `
	DropOffAddress        string    `json:"dropoff_address" `
	DroOffLocation        Location  `json:"dropoff_location" `
	DropOffInstruction    string    `json:"dropoff_instructions" `
	Stage                 string    `json:"stage" `
	ActionIfUndeliverable string    `json:"action_if_undeliverable" `
	PickupAddress         string    `json:"pickup_address" `
	PickUpExternalId      string    `json:"pickup_external_id"`
	PickUpName            string    `json:"pickup_name"`
	PickUpPhone           string    `json:"pickup_phone"`
	PickUpLocation        Location  `json:"pickup_location"`
	PickupTime            time.Time `json:"pickup_time"`
	PickupEstimatedTime   int       `json:"pickup_estimated_time"`
	PickupReferenceTag    string    `json:"pickup_reference_tag" `
	DriverPhone           string    `json:"driver_phone"`
	DriverAllowedVehicles []string  `json:"driver_allowed_vehicles"  `
	DriverExternalId      string    `json:"driver_external_id"`
	Metadata              Metadata  `json:"metadata" `
	CancelReason          string    `json:"cancel_reason"`
}
type Metadata struct {
	CreatedAt time.Time `json:"created_at"`
}
type Location struct {
	Point       string    `json:"point"`
	Coordinates []float64 `json:"coordinates"`
}
type Device struct {
	DeviceId string ` json:"device_id"`
	Kind     string ` json:"kind"` // andriod ,ios
}

func HandleNewOrder(data []byte) {
	var order Order
	err := json.Unmarshal(data, &order)
	if err != nil {
		log.Fatal(err)
	}
	closestDrivers := GetDrivers(order.PickUpLocation.Coordinates)
	var stringids string
	for i, driver := range closestDrivers {
		id := strings.Split(driver.Name, ":")[1]
		if i == len(closestDrivers)-1 {
			stringids = stringids + id
			continue
		}
		stringids = stringids + id + ","
	}
}
func HandleAcceptOrder(data []byte) {}
func HandleDriverAcceptOrder(data []byte)
func HandleMerchRejectOrder(data []byte) {}
func HandleCanceledOrder(data []byte)    {}
func HandleDriverDropOrder(data []byte)  {}
func HandleOrderDelivered(data []byte)   {}
func HandleOrderPickuped(data []byte)    {}
