package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	hub := NewHub()
	go hub.Run()

	http.HandleFunc("/", page)
	http.HandleFunc("/register", register)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		connect(hub, w, r)
	})
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	log.Println(err)
}
