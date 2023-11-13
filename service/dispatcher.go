package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"golang.org/x/exp/slices"
)

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
	req, err := http.Get(fmt.Sprintf(os.Getenv("SERVER_URI")+"/customer/get/%s", id))
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

// driver availability, proximity to the rider, vehicle type or weight or others, and driver ratings // traffic and estimated time of arrival => mongodb
func FindBestDriver(order Order) []Driver {
	var drivers []Driver
	var mindist int = 0
	var maxdist int = 5000 // 1km distance
	req, err := http.Get(fmt.Sprintf(os.Getenv("SERVER_URI")+"/driver/location?lang=%s&lat=%s&mindist=%s&maxdist=%s", order.PickUpLocation.Coordinates[0], order.PickUpLocation.Coordinates[1], mindist, maxdist))
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&drivers)
	if err != nil {
		log.Fatal(err)
	}
	for i, driver := range drivers {
		if !driver.Status {
			drivers = append(drivers[:i], drivers[+1:]...)
		}
		if slices.Contains(order.DriverAllowedVehicles, driver.Vehicle.Type) && driver.Status {
			drivers = append(drivers[:i], drivers[+1:]...)
		}
	}
	// blockedDrivers := tools.RedisClient.LRange(tools.Ctx, "rejects:"+order.Id, 0, -1).Val()
	// for i, driver := range drivers {
	// 	for _, blockeddriver := range blockedDrivers {
	// 		if blockeddriver == driver.Id {
	// 			drivers = append(drivers[:i], drivers[+1:]...)
	// 		}
	// 	}
	// }
	// estimated time
	return drivers
}
