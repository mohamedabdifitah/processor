package utils

import (
	"fmt"
	"math"
)

// this problem was solved by using this https://gist.github.com/hotdang-ca/6c1ee75c48e515aec5bc6db6e3265e49
func CalculateDistance(lat1 float64, lng1 float64, lat2 float64, lng2 float64, unit string) float64 {
	radlat1 := float64(math.Pi * lat1 / 180)
	radlat2 := float64(math.Pi * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(math.Pi * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)
	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515
	if unit == "km" {
		dist = dist * 1.609344
	} else {
		dist = dist * 1.609344 * 1000
	}
	return dist
}
func Converter(v float64) string {
	if v >= 1000 {
		return fmt.Sprintf("%f km", math.Round(v))
	}
	return fmt.Sprintf("%f m", math.Round(v))
}
