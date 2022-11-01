package chat

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[int64]*Client),
		Rooms:      make(map[string]*Room),
		Unregister: make(chan *Client),
		Register:   make(chan RegisterStruct),
		Broadcast:  make(chan Message, 1),
		Request:    make(chan Request),
		EnterRoom:  make(chan JoinRoomStruct),
	}
}

// Core function to run the hub
func (h *Hub) Run() {
	for {
		select {
		// Register a client.
		case client := <-h.Register:
			h.RegisterNewClient(&client.client)
		case request := <-h.Request:
			h.HandleRequest(request)
		// Unregister a client.
		case client := <-h.Unregister:
			h.RemoveClient(client)
		// Broadcast a message to all clients.
		case message := <-h.Broadcast:
			h.HandleMessage(message)
		}
	}
}

// function to handle request
func (h *Hub) HandleRequest(req Request) {
	if req.Event == "message" {
		var msg Message
		jsonBody, err := json.Marshal(req.Data)
		err = json.Unmarshal(jsonBody, &msg)
		if err != nil {
			fmt.Println(err)
		}
		h.HandleMessage(msg)
	}
	if req.Event == "join_room" {
		jrs := &JoinRoomStruct{}
		jsonBody, err := json.Marshal(req.Data)
		if err != nil {
		}
		err = json.Unmarshal(jsonBody, &jrs)
		if err != nil {
			fmt.Println(err)
		}
		h.JoinRoom(h.Clients[int64(jrs.ClientID)], jrs.RoomID)
	}
	if req.Event == "leave_room" {
		jrs := &JoinRoomStruct{}
		jsonBody, err := json.Marshal(req.Data)
		if err != nil {
		}
		err = json.Unmarshal(jsonBody, &jrs)
		if err != nil {
			fmt.Println(err)
		}
		h.LeaveRoom(h.Clients[int64(jrs.ClientID)], jrs.RoomID)
	}

	fmt.Println("Size of clients: ", len(h.Clients))
}

// function to add client
func (h *Hub) RegisterNewClient(client *Client) {
	h.Clients[client.ID] = client
	ClientID := strconv.Itoa(int(client.ID))
	RoomList := h.GetRoomList()
	keys := make([]string, 0, len(RoomList))
	for key := range RoomList {
		keys = append(keys, key)
	}

	list := strings.Join(keys, ",")
	msg := Message{
		Type:      "private",
		Content:   list,
		Recipient: ClientID,
		Sender:    "0",
	}
	h.Broadcast <- msg
	fmt.Println("Size of clients: ", len(h.Clients))
}

// function to join a room
func (h *Hub) JoinRoom(client *Client, roomID string) {
	fmt.Println("Client:", client)
	fmt.Println("RoomID:", roomID)
	fmt.Println("All client:", h.Clients)
	h.CreateRoomIfNotExist(roomID)
	if _, ok := h.Rooms[roomID]; ok {
		if client != nil {
			h.Rooms[roomID].Clients[client.ID] = client
			client.Room[roomID] = h.Rooms[roomID]
		}
	}
}

// function to join a room
func (h *Hub) LeaveRoom(client *Client, roomID string) {
	fmt.Println("Client:", client)
	fmt.Println("RoomID:", roomID)
	fmt.Println("All client:", h.Clients)
	if val, ok := h.Rooms[roomID]; ok {
		if _, ok := val.Clients[client.ID]; ok {
			delete(h.Rooms[roomID].Clients, client.ID)
			fmt.Printf("Room %v, client %v", roomID, h.Rooms[roomID].Clients)
		}
	}
	if _, ok := client.Room[roomID]; ok {
		delete(client.Room, roomID)
	}
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
		close(client.Send)
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
	res := &Response{
		Event:  "notification",
		Status: "success",
		Data:   message,
	}
	room := h.Rooms[RoomID]
	if room != nil {
		fmt.Println("Send notification") // Testing purpose
		h.Rooms[RoomID].Response <- *res
	}
}

// function to handle message based on type of message
func (h *Hub) HandleMessage(message Message) {

	//Check if the message is a type of "message"
	if message.Type == "message" {
		res := &Response{
			Event:  "message",
			Status: "success",
			Data:   message,
		}
		room := h.Rooms[message.Recipient]
		if room != nil {
			room.Response <- *res
		}
	}

	if message.Type == "private" {
		res := &Response{
			Event:  "private",
			Status: "success",
			Data:   message,
		}
		recipientID, _ := strconv.ParseInt(message.Recipient, 10, 64)
		if c, ok := h.Clients[recipientID]; ok {
			c.Send <- *res
		}
	}
}

func (h *Hub) GetRoomList() map[string]*Room {
	return h.Rooms
}
