package main

import (
	"github.com/gorilla/websocket"
)

type Hub struct {
	broadcast  chan WSMessage
	register   chan *User
	unregister chan *User
	clients    map[*User]bool
	rooms      map[string]*Room
}

type Room struct {
	Name      string
	broadcast chan WSMessage
	messages  []WSMessage
	clients   map[*User]bool
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan WSMessage),
		register:   make(chan *User),
		unregister: make(chan *User),
		clients:    make(map[*User]bool),
		rooms:      make(map[string]*Room),
	}
}

func (r *Room) BroadcastRoom(html string) {
	for client := range r.clients {
		client.conn.WriteMessage(websocket.TextMessage, []byte(html))
	}

}

func (r *Room) ConnectUser(u *User) {
	r.clients[u] = true
}

func (r *Room) DisconnectUser(u *User) {
	delete(r.clients, u)
}

func (h *Hub) GetRoom(room string) *Room {
	if h.rooms[room] == nil {
		h.rooms[room] = &Room{
			Name:      room,
			broadcast: make(chan WSMessage),
			messages:  make([]WSMessage, 0),
			clients:   make(map[*User]bool),
		}
	}
	return h.rooms[room]
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			delete(h.clients, client)
			close(client.send)
			client.room.DisconnectUser(client)
		case msg := <-h.broadcast:
			for _, room := range h.rooms {
				room.messages = append(room.messages, msg)
			}
			for client := range h.clients {
				select {
				case client.send <- msg:
				default:
					delete(h.clients, client)
					close(client.send)
					client.room.DisconnectUser(client)
				}
			}
		}

	}
}
