package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"game-server/ws"
	"net/http"
	"os"
	"os/signal"
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

	messages := make(chan interface{})
	go listenForMessages(conn, messages)

	commands := make(chan ws.Command)
	go sendCommands(commands)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		select {
		case message := <-messages:
			fmt.Printf("Received message: %s\n", message)
			go sendCommands(commands)
		case command := <-commands:
			fmt.Printf("Sending command: %s\n", command)
			if err := conn.WriteJSON(command); err != nil {
				fmt.Printf("Error %s sending command %s:", err, command.Cmd)
				return
			}
		case <-interrupt:
			fmt.Println("Interrupt signal received. Exiting...")
			return
		}
	}
}

func listenForMessages(conn *websocket.Conn, messages chan<- interface{}) {
	for {

		var response interface{}
		if err := conn.ReadJSON(&response); err != nil {
			fmt.Println("Error reading server response:", err)
			return
		}
		messages <- response
	}
}

func sendCommands(commands chan<- ws.Command) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter a command (or 'exit' to quit): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			fmt.Println("Exiting the program.")
			return
		}

		var cmd ws.Command
		if err := json.Unmarshal([]byte(input), &cmd); err != nil {
			fmt.Println("Invalid JSON input. Try again.")
			return
		}

		commands <- cmd
	}
}

func init() {
	// Customize the dialer with additional settings if needed.
	websocket.DefaultDialer = &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
	}
}
