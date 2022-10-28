package chat

import (
	"fmt"
)

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[int64]*Client),
		Rooms:      make(map[string]*Room),
		Unregister: make(chan *Client),
		Register:   make(chan RegisterStruct),
		Broadcast:  make(chan Message),
	}
}

// Core function to run the hub
func (h *Hub) Run() {
	for {
		select {
		// Register a client.
		case client := <-h.Register:
			h.RegisterNewClient(&client.client, client.roomID)
			// Unregister a client.
		case client := <-h.Unregister:
			h.RemoveClient(client)
			// Broadcast a message to all clients.
		case message := <-h.Broadcast:
			//Check if the message is a type of "message"
			h.HandleMessage(message)
		}
	}
}

// function to add client
func (h *Hub) RegisterNewClient(client *Client, roomID string) {
	fmt.Println(client.ID)
	h.Clients[client.ID] = client
	h.CreateRoomIfNotExist(roomID)
	h.Rooms[roomID].Clients[client.ID] = client
	fmt.Println("Size of clients: ", len(h.Clients))
}

// function check if room exists and if not create it
func (h *Hub) CreateRoomIfNotExist(RoomID string) {
	if _, ok := h.Rooms[RoomID]; !ok {
		h.Rooms[RoomID] = CreateNewRoom(RoomID)
	}
}

// function to remvoe client from room
func (h *Hub) RemoveClient(client *Client) {
	if _, ok := h.Clients[client.ID]; ok {
		for _, r := range client.Room {
			delete(r.Clients, client.ID) // Delete from all rooms
		}
		delete(h.Clients, client.ID) // Delete from hub
		close(client.send)
		fmt.Println("Removed client")
	}
}

func (h *Hub) SendNotification(RoomID string, content string) {
	message := &Message{
		ID:        "",
		Type:      "notification",
		Sender:    "0",
		Recipient: RoomID,
		Content:   content,
	}
	room := h.Rooms[RoomID]
	if room != nil {
		fmt.Println("Send notification")
		h.Rooms[RoomID].Message <- *message
	}
}

// function to handle message based on type of message
func (h *Hub) HandleMessage(message Message) {

	//Check if the message is a type of "message"
	if message.Type == "message" {
		room := h.Rooms[message.Recipient]
		if room != nil {
			room.Message <- message
		}
	}

	if message.Type == "notification" {
		room := h.Rooms[message.Recipient]
		if room != nil {
			room.Message <- message
		}
	}
}
