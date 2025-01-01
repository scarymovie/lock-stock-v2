package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"lock-stock-v2/external/domain"
	"lock-stock-v2/external/usecase"
	"lock-stock-v2/middleware"
	"log"
	"net/http"
)

type JoinRoom struct {
	joinRoomUseCase usecase.JoinRoom
	roomFinder      domain.RoomFinder
}

func NewJoinRoom(u usecase.JoinRoom, roomFinder domain.RoomFinder) *JoinRoom {
	return &JoinRoom{joinRoomUseCase: u, roomFinder: roomFinder}
}

// ServeHTTP - метод, реализующий интерфейс handlers.JoinRoom.
func (h *JoinRoom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")
	if roomID == "" {
		http.Error(w, "Room ID is required", http.StatusBadRequest)
		return
	}

	room, err := h.roomFinder.FindById(roomID)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if room == nil {
		fmt.Printf("Room %v does not exist", roomID)
	}

	user, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Println(roomID)
	log.Println(user.GetUserId())

	req := usecase.JoinRoomRequest{
		User: user,
		Room: room,
	}

	err = h.joinRoomUseCase.JoinRoom(req)
	if err != nil {
		http.Error(w, "Failed to join room: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("Player %s joined room %s", req.User.GetUserId(), req.Room.GetRoomUid())
}
