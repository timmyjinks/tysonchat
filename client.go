package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type User struct {
	Username string
	hub      *Hub
	room     *Room
	conn     *websocket.Conn
	send     chan WSMessage
	IsAdmin  bool
}

func (u *User) read() {
	defer func() {
		u.hub.unregister <- u
		u.conn.Close()
	}()

	for _, msg := range u.room.messages {
		u.send <- msg
	}

	for {
		var msg Message
		err := u.conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			break
		}
		wsmsg := WSMessage{u.Username, msg.Text, time.Now().Format("2006-01-02 15:04")}
		u.room.messages = append(u.room.messages, wsmsg)

		html := fmt.Sprintf(
			`<div hx-swap-oob="beforeend" id="chat_room"><p>%v %s: %s</p></div>`,
			wsmsg.Sent, wsmsg.Username, wsmsg.Text,
		)

		u.room.BroadcastRoom(html)

		if u.IsAdmin {
			u.hub.broadcast <- wsmsg
		}
	}
}

func (u *User) write() {
	for msg := range u.send {
		html := fmt.Sprintf(
			`<div hx-swap-oob="beforeend" id="chat_room"><p>%v %s: %s</p></div>`,
			msg.Sent, msg.Username, msg.Text,
		)
		err := u.conn.WriteMessage(websocket.TextMessage, []byte(html))
		if err != nil {
			break
		}
	}

}
