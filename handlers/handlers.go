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

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	data := models.HomePage{
		NotEmpty: false,
		Conversations: []models.Conversation{
			{Id: 0, Sender: "Me", Receiver: "Bob"},
			{Id: 1, Sender: "Me", Receiver: "Fred"},
			{Id: 2, Sender: "Me", Receiver: "Joe"},
		},
		Content: nil,
	}
	indexTmpl.Execute(w, data)
}

func OpenConvoHandler(w http.ResponseWriter, r *http.Request) {
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
    } // end dummy data

	// var convo Conversation
	// json.NewDecoder(r.Body).Decode(&convo)
	data := models.HomePage{
		NotEmpty: true,
		Conversations: []models.Conversation{
			{Id: 0, Sender: "Me", Receiver: "Bob"},
			{Id: 1, Sender: "Me", Receiver: "Fred"},
			{Id: 2, Sender: "Me", Receiver: "Joe"},
		},
		Content: texts,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Write([]byte("405 Method not allowed."))
	}
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		session, _ := store.Get(r, "cookie-name")
		session.Values["authenticated"] = true
		session.Save(r, w)
		// redirect to index
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}
	loginTmpl.Execute(w, nil)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// create new user and redirect to index
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
		collection := client.Database("chat").Collection("users")
		result, err := collection.InsertOne(ctx, bsonUser)
		fmt.Println(result)
		database.Close(client, ctx, cancel)

		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}
	signupTmpl.Execute(w, nil)
}