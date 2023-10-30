package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type RegistrationRequest struct {
	Nickname string `json:"nickname"`
}

type RegistrationResponse struct {
	UserID string `json:"userID"`
}

type Command struct {
	Cmd    string `json:"cmd"`
	UserID string `json:"id"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func generateUserID() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func handleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {
		// Define your "join" command data
		cmd := Command{}
		if err := conn.ReadJSON(&cmd); err != nil {
			fmt.Println(err)
			return
		}

		switch cmd.Cmd {
		case "join":
			handleJoinRequest(&cmd, conn)
		default:
			fmt.Printf("Bilinmeyen komut: %s\n", cmd)
		}
	}
}

var waitingRequests = make(map[string]*websocket.Conn)
var rooms = make(map[string][]*websocket.Conn)

func handleJoinRequest(cmd *Command, conn *websocket.Conn) {

	waitingRequests[cmd.UserID] = conn

	for _, v := range waitingRequests {
		response := struct {
			Command string `json:"cmd"`
			Reply   string `json:"reply"`
		}{
			Command: "join",
			Reply:   "Waiting",
		}
		v.WriteJSON(response)
	}
	delete(waitingRequests, cmd.UserID)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var request RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Hatalı istek verisi: %v", err)
		return
	}

	userID := generateUserID()

	response := RegistrationResponse{UserID: userID}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "An error occured: %v", err)
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ws", handleWebSocketConnection)
	r.HandleFunc("/register", register)
	r.HandleFunc("/get", get)

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
