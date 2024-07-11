package main

import (
	"net/http"
    "html/template"
    "fmt"
    "time"
    "encoding/json"
    
    "github.com/gorilla/mux"
)

type Message struct {
    Id int64 `json:"id"`
    Sender string `json:"sender"`
    Receiver string `json:"receiver"`
    Content string `json:"content"`
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

func main() {
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
    })
	
	http.ListenAndServe(":80", router)
}
