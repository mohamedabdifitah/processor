package notification

import (
	"bytes"
	"log"
	"net/http"
)

func SendWebhook(url string, body []byte) bool {
	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		return false
	}
	return true
}
