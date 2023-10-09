package service

import (
	"fmt"

	"github.com/mohamedabdifitah/processor/pubsub"
	"github.com/redis/go-redis/v9"
)

func SearchDrivers(limit int, lang, lat, r float64, unit string, withdist bool) []redis.GeoLocation {
	value, err := pubsub.RedisClient.GeoSearchLocation(pubsub.Ctx, "driver", &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Longitude:  lang,
			Latitude:   lat,
			Radius:     r,
			RadiusUnit: unit,
		},
		WithDist: withdist,
	}).Result()
	if err == redis.Nil {
		fmt.Println("key does not exist")
	} else if err != nil {
		panic(err)
	}
	return value
}

// set driver location using redis
func SetDriverLocation(name string, Longitude float64, Latitude float64) (int64, error) {
	res, err := pubsub.RedisClient.GeoAdd(pubsub.Ctx, "driver", &redis.GeoLocation{
		Name:      name,
		Longitude: Longitude,
		Latitude:  Latitude,
	}).Result()
	if err != nil {
		return 0, err
	}
	return res, nil
}

// driver availability, proximity to the rider, vehicle type or weight or others, and driver ratings // traffic and estimated time of arrival => mongodb
func GetDrivers(coordinates []float64) []redis.GeoLocation {
	closestDrivers := SearchDrivers(5, coordinates[0], coordinates[1], 400, "m", true)
	if len(closestDrivers) == 0 {
		closestDrivers = SearchDrivers(5, coordinates[0], coordinates[1], 3000, "m", true)
	}
	return closestDrivers
}
func MissingDrivers(data []byte) {
}
