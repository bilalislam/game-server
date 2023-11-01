package ws

import (
	"fmt"
	handler "game-server/http"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type Command struct {
	Cmd    string `json:"cmd"`
	UserID string `json:"id"`
}

type ReplyCommand struct {
	Cmd   string `json:"cmd"`
	Reply string `json:"reply"`
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
			fmt.Printf("Bilinmeyen komut: %s\n", cmd)
		}
	}
}

var waitingRequests = make(map[string]*websocket.Conn)
var mu sync.RWMutex

func handleJoinRequest(cmd *Command, conn *websocket.Conn) {
	mu.Lock()
	waitingRequests[cmd.UserID] = conn
	mu.Unlock()

	for _, v := range waitingRequests {
		var reply ReplyCommand
		if handler.Users[cmd.UserID] == nil {
			reply = ReplyCommand{
				Cmd:   cmd.Cmd,
				Reply: "notRegistered",
			}
		} else {
			reply = ReplyCommand{
				Cmd:   cmd.Cmd,
				Reply: "waiting",
			}
		}

		err := v.WriteJSON(reply)
		if err != nil {
			return
		}
	}

	mu.Lock()
	delete(waitingRequests, cmd.UserID)
	mu.Unlock()
}
