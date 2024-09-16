package database

import (
	"fmt"
	"context"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"chat/models"
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
} // end database helper functions

// A function to get a users name given their email address
func RetrieveName(email string) string {
	client, ctx, cancel, err := Connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	defer Close(client, ctx, cancel)
	collection := client.Database("chat").Collection("users")
	
	// Find and decode data
	var user models.User
	collection.FindOne(ctx, bson.D{{"email", email}}).Decode(&user)

	return user.FirstName + " " + user.LastName
}

// A function to get a conversation entry from a given conversation id
func GetConversation(id int64) models.Conversation {
	// connect to the database and get the proper collection
	client, ctx, cancel, err := Connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	defer Close(client, ctx, cancel)
	collection := client.Database("chat").Collection("conversations")

	// Find the conversation and decode the data
	var conversation models.Conversation
	collection.FindOne(ctx, bson.D{{"id", id}}).Decode(&conversation)

	return conversation
}

// Database connection test
func TestDatabase() {
  // connect to a local database server
  client, ctx, cancel, err := Connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	pingerr := client.Ping(ctx, readpref.Primary())
	if pingerr != nil {
		panic(err)
	}

  Close(client, ctx, cancel)
	fmt.Println("Connected Successfully") // print a success message
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

// A function to clear a collection
func ClearCollection(collectionName string) {
	// make a connection to the database
	client, ctx, cancel, err := Connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	defer Close(client, ctx, cancel)
	collection := client.Database("chat").Collection(collectionName)
	collection.DeleteMany(ctx, bson.D{{}})
}
