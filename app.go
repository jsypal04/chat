package main

import (
	"net/http"
    "html/template"
    "fmt"
    "io/ioutil"
    "time"
    // "encoding/json"
    
    "github.com/gorilla/mux"
)

type Message struct {
    Id int64
    Content string
}

type Conversation struct {
    Id int64
    Sender string
    Receiver string
}

type HomePage struct {
    NotEmpty bool
    Conversations []Conversation
    Content []Message
}

func main() {
    router := mux.NewRouter()

    // dummy data
    texts := []Message{}
    for i := 0; i < 12; i++ {
        text := Message{
            Id: time.Now().UnixNano(),
            Content: "Hello World: " + string(i),
        }
        texts = append(texts, text)
    } // end dummy data

    fs := http.FileServer(http.Dir("static/"))
    router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

    tmpl := template.Must(template.ParseFiles("templates/index.html"))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            resp, err := ioutil.ReadAll(r.Body)
            if err != nil {
                fmt.Println(err)
            }
            respStr := string(resp)

            message := Message{
                Id: time.Now().UnixNano(),
                Content: respStr,
            }
            fmt.Println(message)
            return
        }

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
        tmpl.Execute(w, data)
    })
	
	http.ListenAndServe(":80", router)
}
