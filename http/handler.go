package http

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
)

type RegistrationRequest struct {
	Nickname string `json:"nickname"`
}

type RegistrationResponse struct {
	UserID string `json:"userID"`
}

type User struct {
	UserID string `json:"userID"`
	Cups   []int  `json:"cups"`
	Rooms  []int  `json:"rooms"`
}

var Users = make(map[string]*User)
var mu sync.RWMutex

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var request RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := fmt.Fprintf(w, "an error occured: %v", err)
		if err != nil {
			return
		}
		return
	}

	userID := generateUserID()
	response := RegistrationResponse{UserID: userID}
	if Users[userID] == nil {
		mu.Lock()
		Users[userID] = &User{UserID: userID}
		mu.Unlock()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := fmt.Fprintf(w, "An error occured: %v", err)
		if err != nil {
			return
		}
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func generateUserID() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
