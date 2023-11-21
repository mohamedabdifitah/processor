package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/meilisearch/meilisearch-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	mongoURL          = "mongodb://localhost:27017"
	mongoDatabaseName = "resys"
	meiliSearchURL    = "http://localhost:7700"
)

func InitSynchronizer() {
	// Connect to MongoDB
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatal(err)
	}
	defer mongoClient.Disconnect(context.Background())
	err = mongoClient.Ping(Ctx, readpref.Primary())
	if err != nil {
		log.Fatal("error while trying to ping mongo", err)
	}
	fmt.Println("mongo connection established")
	// Connect to MeiliSearch
	meiliSearchClient := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   meiliSearchURL,
		APIKey: os.Getenv("MELLI_API_KEY"),
	})
	_, err = meiliSearchClient.Health()
	if err != nil {
		log.Fatalf("Failed to connect to MeiliSearch: %v", err)
	}
	fmt.Println("Melli client initialized")
	// Define the collections to synchronize
	collectionNames := []string{"merchant", "menu"}

	// Start synchronization for each collection
	for _, collectionName := range collectionNames {
		meiliSearchIndex := meiliSearchClient.Index(collectionName)
		go syncCollection(mongoClient, meiliSearchIndex, collectionName)
	}

	// Keep the program running
	select {}
}
func syncCollection(mongoClient *mongo.Client, meiliSearchIndex *meilisearch.Index, collectionName string) {
	collection := mongoClient.Database(mongoDatabaseName).Collection(collectionName)

	// Open a change stream for the specified collection
	changeStream, err := collection.Watch(context.Background(), mongo.Pipeline{})
	if err != nil {
		log.Fatalf("Error opening change stream for collection %s: %v", collectionName, err)
	}
	defer changeStream.Close(context.Background())

	fmt.Printf("Listening for changes in collection %s\n", collectionName)

	// Process change events
	for changeStream.Next(context.Background()) {
		var change bson.M
		changeStream.Decode(&change)
		operation, ok := getNestedField(change, "operationType")
		if !ok {
			continue
		}
		if os.Getenv("APP_ENV") == "development" {
			fmt.Printf("Syncing with MeiliSearch:")
		}
		switch operation {
		case "delete":
			DeleteDocument(meiliSearchIndex, change)
			continue
		case "update":
			UpdateDocument(meiliSearchIndex, change)
			continue
		case "insert":
			CreateDocument(meiliSearchIndex, change)
			continue
		}
	}
}

// getNestedField retrieves a nested field from a BSON document.
func getNestedField(doc bson.M, key string) (interface{}, bool) {
	for k, v := range doc {
		if k == key {
			return v, true
		}

		// Check if the value is another nested document
		if nestedDoc, ok := v.(bson.M); ok {
			if nestedField, found := getNestedField(nestedDoc, key); found {
				return nestedField, true
			}
		}
	}
	return nil, false
}
func UpdateDocument(index *meilisearch.Index, changeEvent bson.M) {
	var document map[string]interface{} = make(map[string]interface{})
	returnId := changeEvent["documentKey"].(bson.M)["_id"]
	var id string = RemoveObjectID(fmt.Sprintf("%s", returnId))
	// TODO: check error if type is invalid
	updateDescription := changeEvent["updateDescription"].(bson.M)
	truncatedArrays := updateDescription["truncatedArrays"].(bson.A)
	updatedFields := updateDescription["updatedFields"].(bson.M)
	removedFields := updateDescription["removedFields"].(bson.A)
	if len(updatedFields) > 0 {
		for key, field := range updatedFields {
			document[key] = field

		}
		document["id"] = id
		_, err := index.UpdateDocuments(document)
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	if len(removedFields) > 0 {
		var a bson.M
		// meillisearch does not support removing fields from the document
		// copy the original document and remove fileds document and delete from meillisearch
		// copy new document to meillisearch
		err := index.GetDocument(id, &meilisearch.DocumentQuery{}, &a)
		if err != nil {
			fmt.Println(err)
		}
		for _, field := range removedFields {
			delete(a, fmt.Sprintf("%s", field))

		}
		_, err = index.DeleteDocument(id)
		if err != nil {
			fmt.Println(err)
		}
		_, err = index.AddDocuments(a)
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	if len(truncatedArrays) > 0 {
		// TODO: turcated arrays change should be updated also in meillisearch
		fmt.Println(truncatedArrays)
	}
}
func CreateDocument(index *meilisearch.Index, changeEvent bson.M) {
	document := changeEvent["fullDocument"].(bson.M)
	var id string = RemoveObjectID(fmt.Sprintf("%s", document["_id"]))
	delete(document, "_id")
	document["id"] = id
	RemoveSensitiveFields(document)
	task, err := index.AddDocuments(document)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(task)
}
func DeleteDocument(index *meilisearch.Index, changeEvent bson.M) {
	returnId := changeEvent["documentKey"].(bson.M)["_id"]
	var id string = RemoveObjectID(fmt.Sprintf("%s", returnId))
	ok, err := index.Delete(id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ok)
}
func RemoveObjectID(objectid string) string {
	var id string
	pattern := `ObjectID\("([0-9a-fA-F]+)"\)`
	re := regexp.MustCompile(pattern)
	// Find the submatches (captures) in the input string
	matches := re.FindStringSubmatch(objectid)
	// Check if there are matches
	if len(matches) == 2 {
		id = matches[1]
		return id
	}
	return ""
}
func RemoveSensitiveFields(data bson.M) {
	var fields []string = []string{"password", "business_email", "business_phone", "metadata", "device"}
	for _, field := range fields {
		delete(data, field)
	}
}
