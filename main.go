package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

const (
	// environment variables
	cosmosDbEndpointEnvVar = "DOCUMENT_DB_URI"
	cosmosDbKeyEnvVar      = "DOCUMENT_DB_PRIMARY_KEY"

	dbName        = "test-db"
	containerName = "test"

	// help text
	idHelp      = "item id number"
	createHelp  = "item item"
	readHelp    = "read item"
	replaceHelp = "replace item"
	deleteHelp  = "delete item"
)

func main() {
	idFlag := flag.String("id", "", idHelp)
	createFlag := flag.String("create", "", createHelp)
	readFlag := flag.Bool("read", false, readHelp)
	replaceFlag := flag.String("replace", "", replaceHelp)
	deleteFlag := flag.Bool("delete", false, deleteHelp)

	flag.Parse()

	if *idFlag != "" && *createFlag != "" {
		create(*idFlag, *createFlag)
		return
	}

	if *idFlag != "" && *readFlag {
		read(*idFlag)
		return
	}

	if *idFlag != "" && *replaceFlag != "" {
		replace(*idFlag, *replaceFlag)
		return
	}

	if *idFlag != "" && *deleteFlag {
		delete(*idFlag)
		return
	}
}

// connects to cosmos
func connect() *azcosmos.Client {
	key := os.Getenv(cosmosDbKeyEnvVar)
	if key == "" {
		log.Fatal("missing environment variable: ", cosmosDbKeyEnvVar)
		return nil
	}
	endpoint := os.Getenv(cosmosDbEndpointEnvVar)
	if endpoint == "" {
		log.Fatal("missing environment variable: ", cosmosDbEndpointEnvVar)
		return nil
	}

	cred, _ := azcosmos.NewKeyCredential(key)
	client, err := azcosmos.NewClientWithKey(endpoint, cred, nil)
	if err != nil {
		log.Fatalf("failed to init client: %v", err)
		return nil
	}
	return client
}

// create item
func create(id string, val string) {
	client := connect()
	if client == nil {
		log.Fatalf("nil client")
		return
	}

	item := item{
		Id:    id,
		Value: val,
	}
	marshalled, err := json.Marshal(&item)
	if err != nil {
		log.Fatal(err)
	}

	container, err := client.NewContainer(dbName, containerName)
	if err != nil {
		log.Fatalf("failed to get container: %v", err)
		return
	}

	// Create an item
	ctx := context.Background()
	itemResponse, err := container.CreateItem(ctx, azcosmos.NewPartitionKeyString(id), marshalled, nil)
	if err != nil {
		log.Fatalf("failed to create item: %v", err)
		return
	}
	fmt.Printf("successfully created item: %+v", itemResponse)
}

// read item
func read(id string) {
	client := connect()
	if client == nil {
		log.Fatalf("nil client")
		return
	}

	container, err := client.NewContainer(dbName, containerName)
	if err != nil {
		log.Fatalf("failed to get container: %v", err)
		return
	}

	ctx := context.Background()
	itemResponse, err := container.ReadItem(ctx, azcosmos.NewPartitionKeyString(id), id, nil)
	if err != nil {
		log.Fatalf("failed to create item: %v", err)
		return
	}

	itemResponseValue := item{}
	err = json.Unmarshal(itemResponse.Value, &itemResponseValue)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("successfully read item: %+v", itemResponseValue)
}

// replace item
func replace(id string, val string) {
	client := connect()
	if client == nil {
		log.Fatalf("nil client")
		return
	}

	item := item{
		Id:    id,
		Value: val,
	}
	marshalled, err := json.Marshal(&item)
	if err != nil {
		log.Fatal(err)
	}

	container, err := client.NewContainer(dbName, containerName)
	if err != nil {
		log.Fatalf("failed to get container: %v", err)
		return
	}

	// replace item
	ctx := context.Background()
	itemResponse, err := container.ReplaceItem(ctx, azcosmos.NewPartitionKeyString(id), id, marshalled, nil)
	if err != nil {
		log.Fatalf("failed to replace item: %v", err)
		return
	}
	fmt.Printf("successfully replaced item: %+v", itemResponse)
}

// delete item
func delete(id string) {
	client := connect()
	if client == nil {
		log.Fatalf("nil client")
		return
	}

	container, err := client.NewContainer(dbName, containerName)
	if err != nil {
		log.Fatalf("failed to get container: %v", err)
		return
	}

	ctx := context.Background()
	itemResponse, err := container.DeleteItem(ctx, azcosmos.NewPartitionKeyString(id), id, nil)
	if err != nil {
		log.Fatalf("failed to create item: %v", err)
		return
	}
	fmt.Printf("successfully deleted item: %+v", itemResponse)
}

type item struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}
