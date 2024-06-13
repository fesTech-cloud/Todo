package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Client *mongo.Client

// DBInstance creates and returns a MongoDB client instance
func DBInstance(dburl string) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dburl))
	if err != nil {
		log.Fatal(err.Error())
	}

	// Pinging the primary to ensure the connection is established
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Mongodb connected")
	return client
}

// OpenCollection opens a collection in the specified database
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("Todos").Collection(collectionName)
}
