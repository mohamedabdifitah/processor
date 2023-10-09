package utils

import (
	"log"
	"os"
	"path/filepath"
)

func BasePath() string {
	var base string
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if filepath.Base(wd) == "processor" {
		return base
	}
	return base
}
