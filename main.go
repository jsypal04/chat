package main

import (
	"net/http"
	"github.com/gorilla/mux"

	"chat/handlers"
	"chat/database"
)

func main()  {
	// test database connection
	database.TestDatabase()

	router := mux.NewRouter()

	fs := http.FileServer(http.Dir("static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	router.HandleFunc("/", handlers.IndexHandler)
	router.HandleFunc("/id/{id}", handlers.OpenConvoHandler)
	router.HandleFunc("/send", handlers.SendMessageHandler)
	router.HandleFunc("/login", handlers.LoginHandler)
	router.HandleFunc("/signup", handlers.SignupHandler)

	http.ListenAndServe(":80", router)
}