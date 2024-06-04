package main

import (
	"net/http"
    "html/template"
)

type Message struct {
    Id int
    Sender string
    Receiver string
    Content string
}

type Conversation struct {
    Id int
    Receiver string
}

type HomePage struct {
    Conversations []Conversation
}

func main() {
    fs := http.FileServer(http.Dir("static/"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    tmpl := template.Must(template.ParseFiles("templates/index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        data := HomePage{
            Conversations: []Conversation{
                {Id: 0, Receiver: "Bob"},
                {Id: 1, Receiver: "Fred"},
                {Id: 2, Receiver: "Joe"},
            },
        }
        tmpl.Execute(w, data)
	})
	
	http.ListenAndServe(":80", nil)
}
