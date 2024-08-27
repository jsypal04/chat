package handlers

import (
	"fmt"
	"time"
	"net/http"
	"html/template"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gorilla/sessions"

	"chat/models"
	"chat/database"
)

var indexTmpl = template.Must(template.ParseFiles("templates/index.html"))
var loginTmpl = template.Must(template.ParseFiles("templates/login.html"))
var signupTmpl = template.Must(template.ParseFiles("templates/signup.html"))

var store = sessions.NewCookieStore([]byte("super-secret"))

/****************************
Non-exported helper functions
****************************/

// A function to check if a user is authenticated
func isAuthenticated(r *http.Request) (bool) {
	session, err := store.Get(r, "user-cookie")
	if err != nil {
		panic(err)
	}
	return session.Values["authenticated"] == true
}

// A function to send an error code to the client
func sendErrorCode(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	w.Write(jsonResp)
}

/*****************************
Handlers for webpage endpoints
*****************************/

// The hanlder for the main page endpoint
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// check if the user is authenticated
	if !isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
	}

	// dummy data to render on the page
	data := models.HomePage{
		NotEmpty: false,
		Conversations: []models.Conversation{
			{Id: 0, Sender: "Me", Receiver: "Bob"},
			{Id: 1, Sender: "Me", Receiver: "Fred"},
			{Id: 2, Sender: "Me", Receiver: "Joe"},
		},
		Content: nil,
	}
	// render the template for the main page
	indexTmpl.Execute(w, data)
}

// The handler for the login page
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// check that the user is not authenticated
	if isAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}

	// checks if the method is post or get
	if r.Method == http.MethodPost {
		// if the method is post create a sessions for the user and redirect to index
		session, err := store.Get(r, "user-cookie")
		if err != nil {
			panic(err)
		}
		session.Values["authenticated"] = true
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}
	// if the method is get render the login page template
	loginTmpl.Execute(w, nil)
}

// The handler for the logout endpoint
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// get the sessions, set authenticated to false, redirect to the login page
	session, err := store.Get(r, "user-cookie")
	if err != nil {
		panic(err)
	}
	session.Values["authenticated"] = false
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
}

// The handler for the signup page
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	// checks that the user is not authenticated
	if isAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}

	// if the method is post
	if r.Method == http.MethodPost {
		// create new user and save in database
		newUser := models.User{
			FirstName: r.FormValue("firstName"),
			LastName: r.FormValue("lastName"),
			Email: r.FormValue("email"),
			Password: r.FormValue("password"),
		}
		fmt.Println(newUser)

		var bsonUser interface{}
		bsonUser = bson.D{
			{"firstName", r.FormValue("firstName")},
			{"lastName", r.FormValue("lastName")},
			{"email", r.FormValue("email")},
			{"password", r.FormValue("password")},
		}

		// make a connection to the database
		client, ctx, cancel, err := database.Connect("mongodb://localhost:27017")
		if err != nil {
			panic(err)
		}
		// get collection
		collection := client.Database("chat").Collection("users")
		result, err := collection.InsertOne(ctx, bsonUser)
		fmt.Println(result)
		if err != nil {
			panic(err)
		}
		// close the database
		database.Close(client, ctx, cancel)

		// redirect to the main page
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}

	// if the method is get, render the signup page template
	signupTmpl.Execute(w, nil)
}

/*************************
Handlers for API endpoints
*************************/

// The handler for the open conversation endpoint
func OpenConvoHandler(w http.ResponseWriter, r *http.Request) {
	// checks that the user is authenticated
	if !isAuthenticated(r) {
		sendErrorCode(w, r, http.StatusUnauthorized, "Unauthorized")
	} else {
		// dummy data
		texts := []models.Message{}
		for i := 0; i < 12; i++ {
			if i % 2 == 0 {
				text := models.Message{
					Id: time.Now().UnixNano(),
					Sender: "Me",
					Receiver: "Bob",
					Content: "Hello World",
				}
				texts = append(texts, text)
				continue
			}
			text := models.Message{
				Id: time.Now().UnixNano(),
				Sender: "Bob",
				Receiver: "Me",
				Content: "Goodbye World",
			}
			texts = append(texts, text)
		}
	
		data := models.HomePage{
			NotEmpty: true,
			Conversations: []models.Conversation{
				{Id: 0, Sender: "Me", Receiver: "Bob"},
				{Id: 1, Sender: "Me", Receiver: "Fred"},
				{Id: 2, Sender: "Me", Receiver: "Joe"},
			},
			Content: texts,
		} // end dummy data

		// send the data to the client encoded as json
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}

// The handler for the send endpoint
func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	// checks that the user is authenticated
	if !isAuthenticated(r) {
		sendErrorCode(w, r, http.StatusUnauthorized, "Unauthorized")
	} else {
		// if the method is not post, reject it
		if r.Method != http.MethodPost {
			sendErrorCode(w, r, http.StatusMethodNotAllowed, "Only POST requests are allowed at this endpoint")
		}

		// decode the message into a message struct
		var message models.Message
		json.NewDecoder(r.Body).Decode(&message)
		fmt.Println(message)
	
		// convert the message to a bson document
		var bsonMessage interface{}
		bsonMessage = bson.D{
			{"id", message.Id},
			{"convo_id", message.ConvoID},
			{"sender", message.Sender},
			{"receiver", message.Receiver},
			{"content", message.Content},
		}
	
		// make a connection to the database
		client, ctx, cancel, err := database.Connect("mongodb://localhost:27017")
		if err != nil {
			panic(err)
		}
	
		// get the appropriate collection
		collection := client.Database("chat").Collection("messages")
		// insert the message into the collection
		result, err := collection.InsertOne(ctx, bsonMessage)
		if err != nil {
			panic(err)
		}
		fmt.Println(result)
	
		// close the connection to the database
		database.Close(client, ctx, cancel)
	}
}

// The handler for the new conversation endpoint
func NewConversationHandler(w http.ResponseWriter, r *http.Request) {
	// checks that the method is not post
	if r.Method != http.MethodPost {
		sendErrorCode(w, r, http.StatusMethodNotAllowed, "Only POST requests are allowed at this endpoint")
	} else {
		// decode the data into a struct
		var newData models.NewConversationData
		json.NewDecoder(r.Body).Decode(newData)
		fmt.Println(newData)
	}
}