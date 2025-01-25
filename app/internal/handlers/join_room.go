package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"lock-stock-v2/external/domain"
	"lock-stock-v2/external/usecase"
	httpResponse "lock-stock-v2/internal/handlers/response"
	"lock-stock-v2/middleware"
	"log"
	"net/http"
)

type JoinRoom struct {
	joinRoomUseCase     usecase.JoinRoom
	roomFinder          domain.RoomFinder
	roomUsersRepository domain.RoomUserRepository
}

func NewJoinRoom(u usecase.JoinRoom, roomFinder domain.RoomFinder, RoomUserRepository domain.RoomUserRepository) *JoinRoom {
	return &JoinRoom{joinRoomUseCase: u, roomFinder: roomFinder, roomUsersRepository: RoomUserRepository}
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
		log.Println("FindById error:", err)
		http.Error(w, "Failed to find room", http.StatusNotFound)
		return
	}

	if room == nil || room.GetRoomId() == 0 {
		log.Printf("Room %v does not exist or invalid\n", roomID)
		http.Error(w, "Room does not exist or invalid", http.StatusNotFound)
		return
	}

	log.Printf("Room retrieved: ID=%d, UID=%s", room.GetRoomId(), room.GetRoomUid())

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

	log.Printf("roomId before finder %d", room.GetRoomId())
	log.Printf("Handler Room type: %T, value: %+v", room, room)

	roomUsers, err := h.roomUsersRepository.FindByRoom(room)

	if err != nil {
		log.Println("FindByRoom error:", err)
		http.Error(w, "No users found in room", http.StatusOK)
		return
	}

	var response []httpResponse.RoomUserResponse
	for _, ru := range roomUsers {
		respItem := httpResponse.RoomUserResponse{
			RoomUid: ru.GetRoom().GetRoomUid(),
			UserUid: ru.GetUser().GetUserUid(),
		}
		response = append(response, respItem)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Failed to encode roomUsers:", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	log.Printf("Player %s joined room %s", req.User.GetUserUid(), req.Room.GetRoomUid())
}
