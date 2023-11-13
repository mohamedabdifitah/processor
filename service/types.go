package service

import "time"

type Merchant struct {
	Id            string   `json:"id"`
	BusinessName  string   `json:"business_name"`
	BusinessPhone string   `json:"business_phone"`
	Location      Location `json:"location"`
	Address       string   `json:"address"`
	Metadata      struct {
		WebhookEndpoint string `json:"webhook_endpoint"`
	} `json:"metadata"`
	Device Device `json:"-" bson:"device"`
	Closed bool   `json:"closed"`
}
type Driver struct {
	Id        string    `json:"id"`
	GivenName string    `json:"given_name"`
	Age       time.Time `json:"age"`
	Metadata  struct {
		WebhookEndpoint string `json:"webhook_endpoint"`
	} `json:"metadata"`
	Profile  string   `json:"profile"`
	Device   Device   `json:"device"`
	Location Location `json:"location"`
	Status   bool     `json:"status"`
	Vehicle  struct {
		Model   string  `json:"model"`
		Type    string  `json:"type" `
		Payload float64 `json:"payload"`
	} `json:"vehicle"`
	Phone string `json:"phone"`
}
type Customer struct {
	Id         string `json:"id,omitempty"`
	Email      string `json:"email"`
	FamilyName string `json:"family_name"`
	GivenName  string `json:"given_name"`
	Address    string `json:"address"`
	Metadata   struct {
		WebhookEndpoint string `json:"webhook_endpoint"`
	} `json:"metadata"`
	Profile string `json:"profile"`
	Device  Device `json:"device"`
	Phone   string `json:"phone"`
}
type Item struct {
	ItemExternalId string `json:"item_external_id"`
	Quantity       uint   `json:"quantity" `
	Price          uint   `json:"price" `
}
type Order struct {
	Id                    string    `json:"id,omitempty"`
	OrderValue            uint      `json:"order_value"  `
	Type                  string    `json:"type" `
	Items                 []Item    `json:"items"`
	DropOffPhone          string    `json:"dropoff_phone" `
	DropOffExteranlId     string    `json:"dropoff_external_id" `
	DropOffContactName    string    `json:"dropoff_contact_name" `
	DropOffTimeEstimated  int       `json:"dropoff_time_estimated" `
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
	DisplayId             string    `json:"display_id"`
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
