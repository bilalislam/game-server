package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"game-server/ws"
	"net/http"
	"os"
	"strings"
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

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter a command (or 'exit' to quit): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			fmt.Println("Exiting the program.")
			break
		}

		var cmd ws.Command
		if err := json.Unmarshal([]byte(input), &cmd); err != nil {
			fmt.Println("Invalid JSON input. Try again.")
			continue
		}

		if err := conn.WriteJSON(cmd); err != nil {
			fmt.Printf("Error %s sending command %s:", err, cmd.Cmd)
			return
		}

		// Read the server's response, if any
		var response interface{}
		if err := conn.ReadJSON(&response); err != nil {
			fmt.Println("Error reading server response:", err)
			return
		}

		fmt.Printf("Server response: %+v\n", response)
	}
}

func init() {
	// Customize the dialer with additional settings if needed.
	websocket.DefaultDialer = &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
	}
}
