package main

import (
	"github.com/gorilla/websocket"
	"html/template"
	"net/http"
)

var upgrader websocket.Upgrader = websocket.Upgrader{}

type Message struct {
	Text string `json:"message"`
}

type WSMessage struct {
	Username string `json:"username"`
	Text     string `json:"message"`
	Sent     string
}

var tmpl *template.Template = template.Must(template.ParseGlob("./static/*.html"))

func connect(h *Hub, w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	roomName := r.URL.Query().Get("room")
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		return
	}

	room := h.GetRoom(roomName)

	user := &User{Username: username, hub: h, room: room, conn: conn, send: make(chan WSMessage), IsAdmin: false}
	room.ConnectUser(user)
	h.register <- user

	go user.read()
	go user.write()
}

func page(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./static/index.html", "./static/register.html"))
	tmpl.ExecuteTemplate(w, "register.html", nil)
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("username")
	room := r.FormValue("room")
	tmpl := template.Must(template.ParseFiles("./static/index.html", "./static/chat.html"))

	data := map[string]string{"Username": name, "Room": room}

	tmpl.ExecuteTemplate(w, "index.html", data)
}
