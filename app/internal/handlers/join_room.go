package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"lock-stock-v2/external/usecase"
	"net/http"
)

type JoinRoom struct {
	joinRoomUseCase usecase.JoinRoom
}

func NewJoinRoom(u usecase.JoinRoom) *JoinRoom {
	return &JoinRoom{joinRoomUseCase: u}
}

// ServeHTTP - метод, реализующий интерфейс handlers.JoinRoom.
func (h *JoinRoom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	req := usecase.JoinRoomRequest{
		PlayerId: playerId,
		RoomId:   roomID,
	}

	err := h.joinRoomUseCase.JoinRoom(req)
	if err != nil {
		http.Error(w, "Failed to join room: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println(w, "Player %s joined room %s", req.PlayerId, req.RoomId)
}
