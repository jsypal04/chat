package handlers

import (
	"fmt"
	"strconv"
	"net/http"
	"html/template"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gorilla/sessions"
	"github.com/gorilla/mux"

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
		return
	}

	// get this users email
	session, err := store.Get(r, "user-cookie")
	if err != nil {
		panic(err)
	}
	userEmail := session.Values["user"].(string)

	// get this users conversations from the database
	client, ctx, cancel, err := database.Connect("mongodb://localhost:27017")
	if err != nil { 
		panic(err) 
	}
	collection := client.Database("chat").Collection("conversations")

	// Get the data from the database
	filter := bson.D{
		{"$or",
			bson.A{
				bson.D{{"user1", userEmail}},
				bson.D{{"user2", userEmail}},
			},
		},
	}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		panic(err)
	}

	// decode the data and close the database 
	var results []models.Conversation
	if err = cursor.All(ctx, &results); err != nil {
		panic(err)
	}
	database.Close(client, ctx, cancel)

	// Get the names of the users
	var convos []models.RenderedConvo = make([]models.RenderedConvo, len(results), cap(results))
	for i := 0; i < len(convos); i++ {
		user1 := database.RetrieveName(results[i].User1)
		user2 := database.RetrieveName(results[i].User2)
		if results[i].User1 == userEmail {
			convos[i] = models.RenderedConvo{
				Id: results[i].Id,
				ReceiverName: user2,
			}
		} else {
			convos[i] = models.RenderedConvo{
				Id: results[i].Id,
				ReceiverName: user1,
			}
		}
	}

	data := models.HomePage{
		NotEmpty: false,
		Conversations: convos,
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
		return
	}

	// checks if the method is post or get
	if r.Method == http.MethodPost {
		// get the entered email and password
		email := r.FormValue("email")
		password := r.FormValue("password")

		// connect to the database
		client, ctx, cancel, err := database.Connect("mongodb://localhost:27017")
		if err != nil {
			panic(err)
		}
		defer database.Close(client, ctx, cancel)
		collection := client.Database("chat").Collection("users")

		// search for the email entered
		var results models.User
		collection.FindOne(ctx, bson.D{{"email", email}}).Decode(&results)
		// check that the user exists
		if len(results.Password) == 0 {
			data := models.LoginPage{
				Display: "block",
				Issue: "There is no account associated with this email",
			}
			loginTmpl.Execute(w, data)
			return
		}

		// check that the password is correct
		if err = bcrypt.CompareHashAndPassword(results.Password, []byte(password)); err != nil {
			data := models.LoginPage{
				Display: "block",
				Issue: "Password is incorrect",
			}
			loginTmpl.Execute(w, data)
			return
		}

		// if the method is post create a sessions for the user and redirect to index
		session, err := store.Get(r, "user-cookie")
		if err != nil {
			panic(err)
		}
		session.Values["authenticated"] = true
		session.Values["user"] = r.FormValue("email")
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	} else {
		// if the method is get render the login page template
		data := models.LoginPage{
			Display: "none",
			Issue: "",
		}
		loginTmpl.Execute(w, data)
	}
}

// The handler for the logout endpoint
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// this authentication may not be necessary but I decided to do it anyway
	if !isAuthenticated(r) {
		sendErrorCode(w, r, http.StatusForbidden, "You must be signed in to log out")
		return
	}

	// get the sessions, set authenticated to false, redirect to the login page
	session, err := store.Get(r, "user-cookie")
	if err != nil {
		panic(err)
	}
	session.Values["authenticated"] = false
	session.Values["user"] = ""
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
}

// The handler for the signup page
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	// checks that the user is not authenticated
	if isAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	// if the method is post
	if r.Method == http.MethodPost {
		// check that the email does not already exist
		// make a connection to the database
		client, ctx, cancel, err := database.Connect("mongodb://localhost:27017")
		if err != nil {
			panic(err)
		}
		// get collection
		collection := client.Database("chat").Collection("users")

		// search for the email that was entered
		var users []models.User
		cursor, _ := collection.Find(ctx, bson.D{{"email", r.FormValue("email")}})
		if err = cursor.All(ctx, &users); err != nil {
			panic(err)
		}

		if len(users) != 0 {
			// reject the signup credentials
			signupData := models.SignupPage{
				Display: "block",
				Issue: "This email already has an account open.",
			}

			signupTmpl.Execute(w, signupData)
			return
		}

		// hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), 10)
		if err != nil {
			panic(err)
		}

		// create new user and save in database
		newUser := models.User{
			FirstName: r.FormValue("firstName"),
			LastName: r.FormValue("lastName"),
			Email: r.FormValue("email"),
			Password: hashedPassword,
		}
		fmt.Println(newUser)

		result, err := collection.InsertOne(ctx, newUser)
		fmt.Println(result)
		if err != nil {
			panic(err)
		}
		fmt.Println(result)
		// close the database
		database.Close(client, ctx, cancel)

		// redirect to the main page
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	} else {
		// if the method is get, render the signup page template
		signupData := models.SignupPage{
			Display: "none",
			Issue: "",
		}
		signupTmpl.Execute(w, signupData)
	}
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
		// get the current logged in user's email
		session, err := store.Get(r, "user-cookie")
		if err != nil {
			panic(err)
		}
		userEmail := session.Values["user"]

		// get the conversation id
		vars := mux.Vars(r)
		convoIdStr := vars["id"]
		convoId, err := strconv.ParseInt(convoIdStr, 10, 64)
		if err != nil {
			panic(err)
		}

		// get the messages from the database and close the database
		client, ctx, cancel, err := database.Connect("mongodb://localhost:27017")
		if err != nil {
			panic(err)
		}
		collection := client.Database("chat").Collection("messages")
		cursor, err := collection.Find(ctx, bson.D{{"convoID", convoId}})

		var messages []models.Message
		if err = cursor.All(ctx, &messages); err != nil {
			panic(err)
		}
		database.Close(client, ctx, cancel)

		for i := 0; i < len(messages); i++ {
			if messages[i].Sender != userEmail {
				continue
			}
			messages[i].Sender = "Me"
		}

		fmt.Println(messages)
		// send the data to the client encoded as json
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(messages)
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
		var msgData models.NewMessageData
		json.NewDecoder(r.Body).Decode(&msgData)
		fmt.Println(msgData)
		fmt.Println(msgData.ConvoID)

		// Get the current user's email
		session, _ := store.Get(r, "user-cookie")
		userEmail := session.Values["user"]

		// Get the other user's email
		var otherEmail string
		conversation := database.GetConversation(msgData.ConvoID)
		if userEmail == conversation.User1 {
			otherEmail = conversation.User2
		} else {
			otherEmail = conversation.User1
		}
	
		// convert the NewMessageData to Message
		message := models.Message{
			Id: msgData.Id,
			ConvoID: msgData.ConvoID,
			Sender: userEmail.(string),
			Receiver: otherEmail,
			Content: msgData.Content,
		}
	
		// make a connection to the database
		client, ctx, cancel, err := database.Connect("mongodb://localhost:27017")
		if err != nil {
			panic(err)
		}
	
		// get the appropriate collection
		collection := client.Database("chat").Collection("messages")
		// insert the message into the collection
		result, err := collection.InsertOne(ctx, message)
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
	// check that the client is authenticated
	if !isAuthenticated(r) {
		sendErrorCode(w, r, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// checks that the method is not post
	if r.Method != http.MethodPost {
		sendErrorCode(w, r, http.StatusMethodNotAllowed, "Only POST requests are allowed at this endpoint")
	} else {
		// decode the data into a struct
		var newData models.NewConversationData
		json.NewDecoder(r.Body).Decode(&newData)
		
		// get the email of the user who made the request
		session, _ := store.Get(r, "user-cookie")
		var userEmail string = session.Values["user"].(string)
		
		// put the conversation data in a conversation struct
		conversation := models.Conversation{
			Id: newData.Id,
			User1: userEmail,
			User2: newData.Receiver,
		}

		// add the new conversation to the database
		client, ctx, cancel, err := database.Connect("mongodb://localhost:27017")
		if err != nil {
			panic(err)
		}

		// get the collection
		collection := client.Database("chat").Collection("conversations")
		collection.InsertOne(ctx, conversation)

		database.Close(client, ctx, cancel)

		// refresh the home page
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}
}