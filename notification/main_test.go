package notification

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/mohamedabdifitah/processor/utils"
)

func TestMain(t *testing.M) {
	wd, _ := os.Getwd()
	if os.Getenv("APP_ENV") == "development" {
		err := godotenv.Load(filepath.Join(filepath.Dir(wd) + "/.env.local"))
		if err != nil {
			log.Fatal(err)
		}
	}
	utils.AllTemplates.LoadTemplates(filepath.Join(filepath.Dir(wd)+"/assets/json/template.json"), filepath.Dir(wd))
	t.Run()
}
