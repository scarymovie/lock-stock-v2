package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"lock-stock-v2/external/usecase"
	"lock-stock-v2/middleware"
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

	user, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	req := usecase.JoinRoomRequest{
		PlayerId: user.GetId(),
		RoomId:   roomID,
	}

	err = h.joinRoomUseCase.JoinRoom(req)
	if err != nil {
		http.Error(w, "Failed to join room: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println(w, "Player %s joined room %s", req.PlayerId, req.RoomId)
}
