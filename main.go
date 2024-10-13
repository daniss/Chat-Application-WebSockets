package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
)

type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}


var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*websocket.Conn]bool)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Error upgrading connection: %s\n", err)
		return
	}
	clients[conn] = true
	conn.WriteMessage(websocket.TextMessage, []byte("Welcome to the chat!"))
	defer conn.Close()
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading message: %s\n", err)
			return
		}
		for client := range clients {
			if client != conn {
				err = client.WriteMessage(messageType, p)
				if err != nil {
					fmt.Printf("Error writing message: %s\n", err)
					return
				}
			}
		}
		err = conn.WriteMessage(messageType, p)
		if err != nil {
			fmt.Printf("Error writing message: %s\n", err)
			return
		}
	}
}

func main() {
	http.HandleFunc("/ws", handler)
	err := http.ListenAndServe(":8080", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
