package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/exp/slices"
)

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
	Satus    bool     `json:"status"`
	Vehicle  struct {
		Model   string  `json:"model"`
		Type    string  `json:"type" `
		Payload float64 `json:"payload"`
	} `json:"vehicle"`
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

// driver availability, proximity to the rider, vehicle type or weight or others, and driver ratings // traffic and estimated time of arrival => mongodb
func GetDrivers(order Order) []Driver {
	var drivers []Driver
	var offlineDrivers []Driver
	var SmallPayload []Driver
	var mindist int = 0
	var maxdist int = 1000 // 1km distance
	var OrderDrivers map[string][]Driver = make(map[string][]Driver)
	req, err := http.Get(fmt.Sprintf(os.Getenv("SERVER_URI")+"/driver/location?lang=%s&lat=%s&mindist=%s&maxdist=%s", order.PickUpLocation.Coordinates[0], order.PickUpLocation.Coordinates[1], mindist, maxdist))
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&drivers)
	if err != nil {
		log.Fatal(err)
	}
	for _, driver := range drivers {
		if driver.Satus == false {
			offlineDrivers = append(offlineDrivers, driver)
			// drivers = append(drivers[:i], drivers[+1:]...)
		}
		if slices.Contains(order.DriverAllowedVehicles, driver.Vehicle.Type) && driver.Satus == true {
			SmallPayload = append(SmallPayload, driver)
			// drivers = append(drivers[:i], drivers[+1:]...)
		}
	}
	// estimated time
	OrderDrivers["top"] = drivers                // webhook , notification
	OrderDrivers["small_payload"] = SmallPayload // webhook ,notification
	OrderDrivers["offline"] = offlineDrivers     // sms , notifications
	return drivers
}

func GetInformationMer(id string) Merchant {
	var merchant Merchant
	req, err := http.Get(fmt.Sprintf(os.Getenv("SERVER_URI")+"/merchant/get/%s", id))
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&merchant)
	if err != nil {
		log.Fatal(err)
	}
	return merchant
}
func GetInformationCustomer(id string) Customer {
	var customer Customer
	req, err := http.Get(fmt.Sprintf(os.Getenv("SERVER_URI")+"/merchant/get/%s", id))
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&customer)
	if err != nil {
		log.Fatal(err)
	}
	return customer
}
func GetInformationDriver(id string) Driver {
	var driver Driver
	req, err := http.Get(fmt.Sprintf(os.Getenv("SERVER_URI")+"/driver/get/%s", id))
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&driver)
	if err != nil {
		log.Fatal(err)
	}
	return driver
}
