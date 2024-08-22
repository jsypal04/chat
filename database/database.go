package database

import (
	"fmt"
	"context"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
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