package ws

import (
	"fmt"
	handler "game-server/http"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"sort"
	"sync"
	"time"
)

type Command struct {
	Cmd    string `json:"cmd"`
	UserID string `json:"id"`
	RoomId string `json:"room"`
	Data   int    `json:"data"`
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

var mu sync.RWMutex
var waitingRequests = make(map[string]*websocket.Conn)
var rooms = make(map[string][]*UserRoom)

type UserRoom struct {
	UserId     string `json:"userId"`
	RoomData   int    `json:"roomData"`
	IsAnswered bool
	UserData   int `json:"userData"`
	Conn       *websocket.Conn
}

type ScoreBoard struct {
	Event    string    `json:"event"`
	Secret   int       `json:"secret"`
	Rankings []Ranking `json:"rankings"`
}

type Ranking struct {
	Rank        int    `json:"rank"`
	Player      string `json:"player"`
	Guess       int    `json:"guess"`
	DeltaTrophy int    `json:"deltaTrophy"`
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
		case "guess":
			handleGuessRequest(&cmd)
		default:
			fmt.Printf("Unknown command: %s\n", cmd)
		}
	}
}

func handleGuessRequest(cmd *Command) {
	userRooms := rooms[cmd.RoomId]
	for _, user := range userRooms {
		if user.UserId == cmd.UserID {

			user.IsAnswered = true
			user.UserData = cmd.Data

			reply := ReplyCommand{
				Cmd:   cmd.Cmd,
				Reply: "guessReceived",
			}

			err := user.Conn.WriteJSON(reply)
			if err != nil {
				fmt.Println(err)
			}

			isAllAnswered := allUsersAnswered(userRooms)
			if isAllAnswered {
				gameOver(userRooms)
			}

			break
		}
	}
}

func calculateUserRanking(userRooms []*UserRoom) ScoreBoard {
	sortUserGuesses(userRooms)
	var scoreBoard ScoreBoard
	for i, user := range userRooms {
		cup := 0
		if i == 0 {
			cup = 30
		} else if i == 1 {
			cup = 10
		}

		var rankings []Ranking
		rankings = append(rankings, Ranking{
			Rank:        i + 1,
			Player:      user.UserId,
			Guess:       user.UserData,
			DeltaTrophy: cup,
		})

		scoreBoard.Event = "gameOver"
		scoreBoard.Secret = user.RoomData
		scoreBoard.Rankings = rankings
	}

	return scoreBoard
}

func sortUserGuesses(userRooms []*UserRoom) {
	sort.Slice(userRooms, func(i, j int) bool {
		diffI := abs(userRooms[i].UserData - userRooms[i].RoomData)
		diffJ := abs(userRooms[j].UserData - userRooms[i].RoomData)
		return diffI < diffJ
	})
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func gameOver(userRooms []*UserRoom) {
	scoreBoard := calculateUserRanking(userRooms)
	for _, user := range userRooms {
		err := user.Conn.WriteJSON(scoreBoard)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func allUsersAnswered(userRooms []*UserRoom) bool {
	for _, userRoom := range userRooms {
		if !userRoom.IsAnswered {
			return false
		}
	}
	return true
}

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

/*
TODO 3: scoreboard , gameOver
* 3 olmadan oda acılmayacak
* kalan kullancılar eşleşmeyecek
* user bazlı timeoutlar'da bir cevap olacagından her tahmin sonrası odadaki her user'ın cevap verildiği kontrol edilecek
* her user cevap verdiyse zaten room içinde aynı user timeout'u gecerli olacagından ilgili gameOver tahminlerin bittiği oda için server
tarafından calıstırılır.
her guess'de tüm cevapları check et ve scoreboard hesapla, bitti ise gameOver(rooms) ve tüm client'lar scoreboard goster
*/
func MatchUsers() {
	for {
		time.Sleep(5 * time.Second)

		roomID := generateRoomID()
		roomData := rand.Intn(10-1) + 1

		mu.RLock()

		//TODO 1 3erlik odalar
		//if len(waitingRequests) >= 3 {
		//
		//}

		for userID, conn := range waitingRequests {

			if len(rooms[roomID]) >= 2 {
				mu.Unlock()
				break
			}

			rooms[roomID] = append(rooms[roomID], &UserRoom{
				UserId:   userID,
				RoomData: roomData,
				Conn:     conn,
			})

			delete(waitingRequests, userID)
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
