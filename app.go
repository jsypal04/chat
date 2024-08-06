package main

import (
	"net/http"
    "html/template"
    "context"
    "fmt"
    "time"
    "encoding/json"
    
    "github.com/gorilla/mux"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Message struct {
    Id int64 `json:"id" bson:"id"`
    ConvoID int64 `json:"convoID" bson:"convo_id"`
    Sender string `json:"sender" bson:"sender"`
    Receiver string `json:"receiver" bson:"receiver"`
    Content string `json:"content" bson:"content"`
}

type Conversation struct {
    Id int64 `json:"id"`
    Sender string `json:"sender"`
    Receiver string `json:"receiver"`
}

type HomePage struct {
    NotEmpty bool `json:"notEmpty"`
    Conversations []Conversation `json:"conversations"`
    Content []Message `json:"content"`
}

// Database helper functions
func connect(uri string)(*mongo.Client, context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}

func close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
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
func testDatabase() {
    // connect to a local database server
    client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected Successfully") // print a success message

	defer cancel()
	defer func() {
		err := client.Disconnect(ctx)
		if err != nil {
			panic(err)
		}
	}()

	pingerr := client.Ping(ctx, readpref.Primary())
	if pingerr != nil {
		panic(err)
	}
}

func main() {
    // test database connection
    testDatabase()

    router := mux.NewRouter()

    // dummy data
    texts := []Message{}
    for i := 0; i < 12; i++ {
        if i % 2 == 0 {
            text := Message{
                Id: time.Now().UnixNano(),
                Sender: "Me",
                Receiver: "Bob",
                Content: "Hello World",
            }
            texts = append(texts, text)
            continue
        }
        text := Message{
            Id: time.Now().UnixNano(),
            Sender: "Bob",
            Receiver: "Me",
            Content: "Goodbye World",
        }
        texts = append(texts, text)
    } // end dummy data

    fs := http.FileServer(http.Dir("static/"))
    router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

    tmpl := template.Must(template.ParseFiles("templates/index.html"))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        data := HomePage{
            NotEmpty: false,
            Conversations: []Conversation{
                {Id: 0, Sender: "Me", Receiver: "Bob"},
                {Id: 1, Sender: "Me", Receiver: "Fred"},
                {Id: 2, Sender: "Me", Receiver: "Joe"},
            },
            Content: nil,
        }
        tmpl.Execute(w, data)
	})

    router.HandleFunc("/id/{id}", func(w http.ResponseWriter, r *http.Request) {
        // var convo Conversation
        // json.NewDecoder(r.Body).Decode(&convo)
        data := HomePage{
            NotEmpty: true,
            Conversations: []Conversation{
                {Id: 0, Sender: "Me", Receiver: "Bob"},
                {Id: 1, Sender: "Me", Receiver: "Fred"},
                {Id: 2, Sender: "Me", Receiver: "Joe"},
            },
            Content: texts,
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(data)
    })

    router.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            w.Write([]byte("405 Method not allowed."))
        }
        var message Message
        json.NewDecoder(r.Body).Decode(&message)
        fmt.Println(message)

        // make a connection to the database
        client, ctx, cancel, err := connect("mongodb://localhost:27017")
        if err != nil {
            panic(err)
        }

        // get the appropriate collection
        collection := client.Database("chat").Collection("messages")
        result, err := collection.InsertMany(ctx, message)
        if err != nil {
            panic(err)
        }
        fmt.Println(result)

        close(client, ctx, cancel)
    })
	
	http.ListenAndServe(":80", router)
}
