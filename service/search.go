package service

import (
	"fmt"
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

func InitElasticSearch() {
	cfg := elasticsearch.Config{
		Username: os.Getenv("ELAST_SEARCH_USER"),
		Password: os.Getenv("ELAST_SEARCH_PASS"),
		Addresses: []string{
			"http://localhost:9200",
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatal(err)
	}
	Create(es)
	fmt.Println("connected to elastic search")
	defer res.Body.Close()
	// SearchDocument(es)
}
func Create(es *elasticsearch.Client) {
	doc := struct {
		Title string `json:"title"`
	}{Title: "Test"}

	res, err := es.Index("test_ree", esutil.NewJSONReader(&doc), es.Index.WithRefresh("true"))
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	log.Println(res)
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": "test",
			},
		},
	}

	res, err = es.Search(
		es.Search.WithIndex("test_ree"),
		es.Search.WithBody(esutil.NewJSONReader(&query)),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	log.Println(res)
}
func SearchDocument(es *elasticsearch.Client) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": "test",
			},
		},
	}

	res, err := es.Search(
		es.Search.WithIndex("test_re"),
		es.Search.WithBody(esutil.NewJSONReader(&query)),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	log.Println(res)
}
