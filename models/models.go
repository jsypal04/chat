package models

// struct definitions for users
type User struct {
    FirstName string
    LastName string
    Email string
    Password string
}

// struct definitions for messages and conversations
type Message struct {
    Id int64 `json:"id"`
    ConvoID int64 `json:"convoID"`
    Sender string `json:"sender"`
    Receiver string `json:"receiver"`
    Content string `json:"content"`
}

type Conversation struct {
    Id int64 `json:"id"`
    Sender string `json:"sender"`
    Receiver string `json:"receiver"`
}


// struct definition for home page data
type HomePage struct {
    NotEmpty bool `json:"notEmpty"`
    Conversations []Conversation `json:"conversations"`
    Content []Message `json:"content"`
}