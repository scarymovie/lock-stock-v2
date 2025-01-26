package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"lock-stock-v2/external/domain"
	"lock-stock-v2/external/usecase"
	services "lock-stock-v2/internal/domain/service"
	"lock-stock-v2/internal/handlers/helpers"
	httpResponse "lock-stock-v2/internal/handlers/response"
	"net/http"
)

type JoinRoom struct {
	joinRoomUseCase usecase.JoinRoom
	roomFinder      domain.RoomFinder
	roomUserService *services.RoomUserService
}

func NewJoinRoom(u usecase.JoinRoom, roomFinder domain.RoomFinder, roomUserService *services.RoomUserService) *JoinRoom {
	return &JoinRoom{joinRoomUseCase: u, roomFinder: roomFinder, roomUserService: roomUserService}
}

func (h *JoinRoom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	req := usecase.JoinRoomRequest{User: user, Room: room}
	if err := h.joinRoomUseCase.JoinRoom(req); err != nil {
		http.Error(w, "Failed to join room: "+err.Error(), http.StatusInternalServerError)
		return
	}

	roomUsers, err := h.roomUserService.GetUsersByRoom(room)
	if err != nil {
		http.Error(w, "Failed to get users in room: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var response []httpResponse.RoomUserResponse
	for _, ru := range roomUsers {
		response = append(response, httpResponse.RoomUserResponse{
			RoomUid:  ru.GetRoom().GetRoomUid(),
			UserUid:  ru.GetUser().GetUserUid(),
			UserName: ru.GetUser().GetUserName(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
