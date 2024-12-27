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

func NewJoinRoomHandler() *JoinRoom {
	return &JoinRoom{}
}

// ServeHTTP - метод, который реализует интерфейс JoinRoom.
func (h *JoinRoom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")
	if roomID == "" {
		http.Error(w, "Room ID is required", http.StatusBadRequest)
		return
	}
	playerId := chi.URLParam(r, "playerId")
	if playerId == "" {
		http.Error(w, "player ID is required", http.StatusBadRequest)
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
