package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	serverAddr := "ws://localhost:8080/ws"

	// Establish a WebSocket connection to the server
	conn, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
	if err != nil {
		fmt.Println("WebSocket connection error:", err)
		return
	}
	defer conn.Close()

	// Define your "join" command data
	command := struct {
		Command string `json:"cmd"`
		UserID  string `json:"id"`
	}{
		Command: "join",
		UserID:  "kquKIltb",
	}

	// Send the "join" command to the server
	if err := conn.WriteJSON(command); err != nil {
		fmt.Println("Error sending 'join' command:", err)
		return
	}

	// Read the server's response, if any
	var response interface{}
	if err := conn.ReadJSON(&response); err != nil {
		fmt.Println("Error reading server response:", err)
		return
	}

	fmt.Printf("Server response: %+v\n", response)
	select {}
}

func init() {
	// Customize the dialer with additional settings if needed.
	websocket.DefaultDialer = &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
	}
}
