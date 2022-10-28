package chat

import "github.com/gorilla/websocket"

// Hub is a struct that holds all the clients and the messages that are sent to them
type Hub struct {
	// Registered clients.
	Clients map[int64]*Client // client list, room list
	// Registered clients.
	Rooms map[string]*Room // client list, room list
	//Unregistered clients.
	Unregister chan *Client
	// Register requests from the clients.
	Register chan RegisterStruct
	// Inbound messages from the clients.
	Broadcast chan Message
}

type Room struct {
	// Room ID
	ID string
	// Client in room
	Clients map[int64]*Client
	// Inbound messages from the clients
	Message chan Message
}

// Client struct for websocket connection and message sending
type Client struct {
	ID   int64 // individual id
	Conn *websocket.Conn
	send chan Message
	hub  *Hub
	Room map[string]*Room // List of room client involved
}

// Message struct to hold message data
type Message struct {
	Type      string `json:"type"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
	ID        string `json:"id"` //individual id
}
