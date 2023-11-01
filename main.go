package main

import (
	"fmt"
	handler "game-server/http"
	wsHandler "game-server/ws"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ws", wsHandler.HandleWebSocketConnection)
	r.HandleFunc("/register", handler.Register)
	r.HandleFunc("/get", handler.Get)

	http.Handle("/", r)

	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			panic("Error starting HTTP server: " + err.Error())
		}
	}()

	fmt.Println("HTTP server started on :8080")

	wsServer := &http.Server{Addr: ":8081"}
	wsServer.Handler = r

	go func() {
		err := wsServer.ListenAndServe()
		if err != nil {
			panic("Error starting WebSocket server: " + err.Error())
		}
	}()

	fmt.Println("WebSocket server started on :8081")

	select {}
}
