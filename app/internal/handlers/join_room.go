package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"lock-stock-v2/external/usecase"
	"net/http"
)

type JoinRoom struct {
}

func NewJoinRoom() *JoinRoom {
	return &JoinRoom{}
}

// ServeHTTP - метод, реализующий интерфейс handlers.JoinRoom.
func (h *JoinRoom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("join room")
	roomID := chi.URLParam(r, "roomId")
	if roomID == "" {
		http.Error(w, "Room ID is required", http.StatusBadRequest)
		return
	}
	playerId := chi.URLParam(r, "playerId")
	if playerId == "" {
		http.Error(w, "Player ID is required", http.StatusBadRequest)
		return
	}

	var req usecase.JoinRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Player %s joined room %s", req.PlayerId, roomID)
}
