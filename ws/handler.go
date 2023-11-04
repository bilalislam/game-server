package ws

import (
	"fmt"
	handler "game-server/http"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type Command struct {
	Cmd    string `json:"cmd"`
	UserID string `json:"id"`
}

type ReplyCommand struct {
	Cmd   string `json:"cmd"`
	Reply string `json:"reply"`
}

type ReplyEvent struct {
	Event string `json:"event"`
	Room  string `json:"room"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandleWebSocketConnection(w http.ResponseWriter, r *http.Request) {

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
			fmt.Printf("Unknown command: %s\n", cmd)
		}
	}
}

var waitingRequests = make(map[string]*websocket.Conn)
var rooms = make(map[string][]*websocket.Conn)

var mu sync.RWMutex

func handleJoinRequest(cmd *Command, conn *websocket.Conn) {
	mu.Lock()
	waitingRequests[cmd.UserID] = conn
	mu.Unlock()

	reply := ReplyCommand{
		Cmd: cmd.Cmd,
	}

	if handler.Users[cmd.UserID] == nil {
		reply.Reply = "notRegistered"
		delete(waitingRequests, cmd.UserID)
	} else {
		reply.Reply = "waiting"
	}

	err := conn.WriteJSON(reply)
	if err != nil {
		fmt.Println(err)
	}
}

func MatchUsers() {
	for {
		time.Sleep(5 * time.Second)
		fmt.Print("run...")

		mu.Lock()
		for userID := range waitingRequests {
			fmt.Printf("user %s", userID)
		}
		mu.Unlock()

		roomID := generateRoomID()

		mu.RLock()
		for userID, conn := range waitingRequests {
			if len(rooms[roomID]) >= 2 {
				mu.Unlock()
				break
			}

			rooms[roomID] = append(rooms[roomID], conn)
			delete(waitingRequests, userID)
			fmt.Printf("roomId %s", roomID)
			fmt.Print("deleted...")

			reply := ReplyEvent{
				Event: "joinedRoom",
				Room:  roomID,
			}

			if err := conn.WriteJSON(reply); err != nil {
				fmt.Println(err)
			}
		}
		mu.RUnlock()
	}
}

func generateRoomID() string {
	return fmt.Sprintf("room%d", time.Now().UnixNano())
}
