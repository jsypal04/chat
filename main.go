package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"

	"chat/handlers"
	"chat/database"
)

func printDb(collectionName string) {
	client, ctx, cancel, err := database.Connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	defer database.Close(client, ctx, cancel)
	collection := client.Database("chat").Collection(collectionName)
	cursor, err := collection.Find(ctx, bson.D{{}}, nil)
	if err != nil {
		panic(err)
	}
	var results []bson.M
	err2 := cursor.All(ctx, &results)
	if err2 != nil {
		panic(err2)
	}
	for _, result := range results {
		fmt.Println(result)
	}
}

func main()  {
	// test database connection
	database.TestDatabase()
	printDb("users")

	router := mux.NewRouter()

	fs := http.FileServer(http.Dir("static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Webpage endpoints
	router.HandleFunc("/", handlers.IndexHandler)
	router.HandleFunc("/login", handlers.LoginHandler)
	router.HandleFunc("/logout", handlers.LogoutHandler)
	router.HandleFunc("/signup", handlers.SignupHandler)

	// API endpoints
	router.HandleFunc("/id/{id}", handlers.OpenConvoHandler)
	router.HandleFunc("/send", handlers.SendMessageHandler)
	router.HandleFunc("/new-convo", handlers.NewConversationHandler)

	http.ListenAndServe(":80", router)
}