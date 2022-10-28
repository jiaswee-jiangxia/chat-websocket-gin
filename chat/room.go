package chat

func CreateNewRoom(RoomID string) *Room {
	room := Room{
		ID:      RoomID,
		Clients: make(map[int64]*Client, 0),
		Message: make(chan Message),
	}
	go room.Open()
	return &room
}

func (r *Room) Open() {
	for {
		select {
		case msg := <-r.Message:
			r.Broadcast(msg)
		}
	}
}

func (r *Room) Broadcast(msg Message) {
	for client := range r.Clients {
		select {
		case r.Clients[client].send <- msg:
		default:
			delete(r.Clients, client)
		}
	}
}
