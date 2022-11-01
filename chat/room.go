package chat

func CreateNewRoom(RoomID string) *Room {
	room := Room{
		ID:       RoomID,
		Clients:  make(map[int64]*Client, 0),
		Response: make(chan Response),
	}
	go room.Open()
	return &room
}

func (r *Room) Open() {
	for {
		select {
		case res := <-r.Response:
			r.Broadcast(res)
		}
	}
}

func (r *Room) Broadcast(res Response) {
	for client := range r.Clients {
		select {
		case r.Clients[client].Send <- res:
		default:
			delete(r.Clients, client)
		}
	}
}
