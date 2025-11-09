package main

type Hub struct {
	broadcast  chan WSMessage
	register   chan *User
	unregister chan *User
	clients    map[*User]bool
	messages   map[string][]WSMessage
	rooms      map[string]Room
}

type Room struct {
	Name      string
	broadcast chan Message
}

func NewRoom(name string) *Room {
	return &Room{
		Name:      name,
		broadcast: make(chan Message),
	}

}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan WSMessage),
		register:   make(chan *User),
		unregister: make(chan *User),
		clients:    make(map[*User]bool),
		messages:   make(map[string][]WSMessage, 0),
		rooms:      make(map[string]Room),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			delete(h.clients, client)
			close(client.send)
		case msg := <-h.broadcast:
			h.messages["string"] = append(h.messages["string"], msg)
			for client := range h.clients {
				select {
				case client.send <- msg:
				default:
					delete(h.clients, client)
					close(client.send)
				}
			}
		}

	}
}
