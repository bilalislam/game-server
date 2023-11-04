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

	messages := make(chan ws.ReplyEvent)
	go listenForMessages(conn, messages)

	commands := make(chan ws.Command)
	go sendCommands("default", commands)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		select {
		case message := <-messages:
			if message.Event != "" {
				if jsonResult, err := json.Marshal(message); err == nil {
					fmt.Println(string(jsonResult))
				}
				go sendCommands(message.Event, commands)
			}
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

// TODO 2 : client bazlı guess komutu için 20 sn timeout
// eger timeout varsa reply olarak don ve timeout oldugunu set et,kupa verme ,sıra da -1 olmalı
// aynı zamanda guess verememesi için cli'ı freeze et , sadece dinlesin
func listenForMessages(conn *websocket.Conn, messages chan<- ws.ReplyEvent) {
	for {

		var response ws.ReplyEvent
		if err := conn.ReadJSON(&response); err != nil {
			fmt.Println("Error reading server response:", err)
			return
		}
		messages <- response
	}
}

func sendCommands(commandType string, commands chan<- ws.Command) {
	for {
		var reader = bufio.NewReader(os.Stdin)
		if commandType == "default" {
			fmt.Print("Enter a join command (or 'exit' to quit): ")
		} else if commandType == "joinedRoom" {
			fmt.Print("Enter a guess command (or 'exit' to quit): ")
		}

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
