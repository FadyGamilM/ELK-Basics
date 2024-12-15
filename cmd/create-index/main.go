package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

func main() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	indexName := flag.String("index-name", "test", "")
	indexReplicas := flag.Int("index-replicas", 2, "")
	indexShards := flag.Int("index-shards", 3, "")

	flag.Parse()

	indexSettings := map[string]interface{}{
		"settings": map[string]interface{}{
			"index": map[string]interface{}{
				"number_of_shards":   3,
				"number_of_replicas": 2,
			},
		},
	}

	// Step 3: Convert the settings to JSON
	settingsJSON, err := json.Marshal(indexSettings)
	if err != nil {
		log.Fatalf("Error marshaling index settings: %s", err)
	}

	// Step 4: Send the request to create the index
	res, err := es.Indices.Create(
		*indexName,
		es.Indices.Create.WithBody(bytes.NewReader(settingsJSON)),
		es.Indices.Create.WithContext(context.Background()),
	)
	if err != nil {
		log.Fatalf("Cannot create index: %s", err)
	}
	defer res.Body.Close()

	// Step 5: Log the response
	fmt.Printf("Response: %s\n", res.String())

	// Step 6: Verify index creation
	res, err = es.Indices.Get([]string{indexName})
	if err != nil {
		log.Fatalf("Cannot get index info: %s", err)
	}
	defer res.Body.Close()
	fmt.Printf("Index created successfully: %s\n", res.String())
}
