package database

import (
	"fmt"
	"context"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Database helper functions
func Connect(uri string)(*mongo.Client, context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}

func Close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
    defer cancel()
    defer func() {
        err := client.Disconnect(ctx)
        if err != nil {
            panic(err)
        }
    }()
}

// end database helper functions

// Database connection test
func TestDatabase() {
    // connect to a local database server
    client, ctx, cancel, err := Connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected Successfully") // print a success message

	pingerr := client.Ping(ctx, readpref.Primary())
	if pingerr != nil {
		panic(err)
	}

    Close(client, ctx, cancel)
}

// A function to print out the contents of a collection from the database
func PrintCollection(collectionName string) {
	// make a database connection
	client, ctx, cancel, err := Connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	defer Close(client, ctx, cancel)
	// get the proper collection
	collection := client.Database("chat").Collection(collectionName)
	// find all entries
	cursor, err := collection.Find(ctx, bson.D{{}}, nil)
	if err != nil {
		panic(err)
	}
	// print the results
	var results []bson.M
	err2 := cursor.All(ctx, &results)
	if err2 != nil {
		panic(err2)
	}
	for _, result := range results {
		fmt.Println(result)
	}
}