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
	// Join room requests from the clients.
	EnterRoom chan JoinRoomStruct
	Request   chan Request
	Response  chan Response
}

type Room struct {
	// Room ID
	ID string
	// Client in room
	Clients map[int64]*Client
	// Inbound messages from the clients
	Response chan Response
}

// Client struct for websocket connection and message sending
type Client struct {
	ID       int64 // individual id
	Nickname string
	Conn     *websocket.Conn
	Send     chan Response
	hub      *Hub
	Room     map[string]*Room // List of room client involved
}

// Message struct to hold message data
type Message struct {
	Type      string `json:"type"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
	ID        string `json:"id"` //individual id
}

type Request struct {
	ID     string      `json:"id"`             // Not used yet, placeholder for future need
	Sender string      `json:"sender"`         // Client who send the request, to be replaced by token check
	Event  string      `json:"event"`          // Event type
	Data   interface{} `json:"data,omitempty"` // 数据 json
}

type Response struct {
	ID     string      `json:"id"`     // Not used yet, placeholder for future need
	Status string      `json:"status"` // success or failed
	Event  string      `json:"event"`  // Event type
	Data   interface{} `json:"data"`   // 数据 json
}

type RegisterStruct struct {
	client Client
}

type JoinRoomStruct struct {
	ClientID int    `json:"clientID"`
	RoomID   string `json:"roomID"`
}
