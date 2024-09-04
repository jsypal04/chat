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

	// Webpage endpoints
	router.HandleFunc("/", handlers.IndexHandler)
	router.HandleFunc("/login", handlers.LoginHandler)
	router.HandleFunc("/logout", handlers.LogoutHandler)
	router.HandleFunc("/signup", handlers.SignupHandler)

	// API endpoints
	router.HandleFunc("/id/{id}", handlers.OpenConvoHandler)
	router.HandleFunc("/send", handlers.SendMessageHandler)
	router.HandleFunc("/new-convo", handlers.NewConversationHandler)
	router.HandleFunc("/get-users", handlers.GetUsersHandler)

	http.ListenAndServe(":80", router)
}