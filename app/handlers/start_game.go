package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"lock-stock-v2/external/domain"
	"lock-stock-v2/external/usecase"
	"lock-stock-v2/handlers/helpers"
	services "lock-stock-v2/internal/domain/service"
	"net/http"
)

type StartGame struct {
	roomFinder       domain.RoomFinder
	roomUserService  *services.RoomUserService
	startGameUsecase usecase.StartGame
}

func NewStartGame(roomFinder domain.RoomFinder, roomUserService *services.RoomUserService, startGameUsecase usecase.StartGame) *StartGame {
	return &StartGame{roomFinder: roomFinder, roomUserService: roomUserService, startGameUsecase: startGameUsecase}
}

func (h *StartGame) Do(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomId")
	room, err := helpers.GetRoomById(h.roomFinder, roomID)
	if err != nil {
		http.Error(w, err.Error(), err.(*helpers.RoomNotFoundError).Code)
		return
	}

	user, err := helpers.GetUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), err.(*helpers.UserNotFoundError).Code)
		return
	}
	roomUsers, err := h.roomUserService.GetUsersByRoom(room)
	if err != nil {
		http.Error(w, "Failed to get users in room: "+err.Error(), http.StatusInternalServerError)
		return
	}

	isUserInRoom := h.roomUserService.IsUserInRoom(roomUsers, user)

	if !isUserInRoom {
		http.Error(w, "User is not in the room", http.StatusForbidden)
		return
	}

	req := usecase.StartGameRequest{
		Room: room,
		User: user,
	}
	if err := h.startGameUsecase.StartGame(req); err != nil {
		http.Error(w, "Failed to start game: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Game started"})
}
