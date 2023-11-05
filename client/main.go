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
	serverAddr := "ws://localhost:8081/ws"

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
	go sendCommands("default", commands)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	var timeout time.Timer
	var isGuessCommandSent bool
	for {
		select {
		case message := <-messages:
			var cmd ws.ReplyCommand
			var event ws.ReplyEvent
			jsonResult, _ := json.Marshal(message)

			if err := json.Unmarshal(jsonResult, &cmd); err != nil {
				fmt.Println("Invalid JSON input. Try again.")
				return
			}

			if err = json.Unmarshal(jsonResult, &event); err != nil {
				fmt.Println("Invalid JSON input. Try again.")
				return
			}

			fmt.Println(string(jsonResult))

			if event.Event != "" {

				if event.Event == "joinedRoom" {
					timeout = *time.NewTimer(20 * time.Second)
				}

				go sendCommands(event.Event, commands)
			} else if cmd.Reply != "" {
				go sendCommands(cmd.Reply, commands)
			}
		case command := <-commands:
			fmt.Printf("Sending command: %s\n", command)
			if err := conn.WriteJSON(command); err != nil {
				fmt.Printf("Error %s sending command %s:", err, command.Cmd)
				return
			}

			if command.Cmd == "guess" {
				isGuessCommandSent = true
				fmt.Printf("please wait game result ... \n")
			}
		case <-timeout.C:
			if !isGuessCommandSent {
				fmt.Println("Timeout occurred. No command received in 20 seconds.")
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

func sendCommands(commandType string, commands chan<- ws.Command) {
	for {
		var reader = bufio.NewReader(os.Stdin)
		if commandType == "default" {
			fmt.Println("Enter a join command (or 'exit' to quit): ")
		} else if commandType == "waiting" {
			fmt.Println("please wait ...")
		} else if commandType == "joinedRoom" {
			fmt.Println("Enter a guess command in 20 seconds (or 'exit' to quit): ")
		} else if commandType == "notRegistered" {
			fmt.Println("Enter a join command (or 'exit' to quit): ")
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
