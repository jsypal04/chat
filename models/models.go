package models

// struct definitions for users
type User struct {
    FirstName string `bson:"firstName"`
    LastName string `bson:"lastName"`
    Email string `bson:"email"`
    Password string `bson:"password"`
}

// struct definitions for messages and conversations
type Message struct {
    Id int64 `bson:"id" json:"id"`
    ConvoID int64 `bson:"convoID" json:"convoID"`
    Sender string `bson:"sender" json:"sender"`
    Receiver string `bson:"receiver" json:"receiver"`
    Content string `bson:"content" json:"content"`
}

type NewMessageData struct {
    Id int64 `json:"id"`
    ConvoID int64 `json:"convoID"`
    Content string `json:"content"`
}

type Conversation struct {
    Id int64 `bson:"id"`
    User1 string `bson:"user1"`
    User2 string `bson:"user2"`
}

type NewConversationData struct {
    Id int64 `json:"id"`
    Receiver string `json:"receiver"`
}

type RenderedConvo struct {
    Id int64 `json:"id"`
    ReceiverName string `json:"receiverName"`
}

// struct definition for home page data
type HomePage struct {
    NotEmpty bool `json:"notEmpty"`
    Conversations []RenderedConvo `json:"conversations"`
    Content []Message `json:"content"`
}