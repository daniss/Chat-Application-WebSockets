package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Error upgrading connection: %s\n", err)
		return
	}
	defer conn.Close()
	for {
		messageType, p, err := conn.ReadMessage()
		fmt.Printf("Received message: %s\n", p)
		if err != nil {
			fmt.Printf("Error reading message: %s\n", err)
			return
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
