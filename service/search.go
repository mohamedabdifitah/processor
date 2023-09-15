package service

import (
	"fmt"
	"log"

	"github.com/mohamedabdifitah/processor/db"
	"github.com/redis/go-redis/v9"
)

// var Ctx = context.Background()

func SearchDrivers(limit int, lat, lang, r float64, unit string, withdist bool) []redis.GeoLocation {
	value, err := db.RedisClient.GeoSearchLocation(db.Ctx, "driver", &redis.GeoSearchLocationQuery{
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

// publish topic to redis channel
func PublishTopic(topic string, message interface{}) {
	err := db.RedisClient.Publish(db.Ctx, topic, message).Err()
	if err != nil {
		log.Fatal(err)
	}
}
