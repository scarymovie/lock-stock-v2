package room

import (
	"encoding/json"
	"lock-stock-v2/handlers/http/helpers"
	"lock-stock-v2/internal/domain/room/repository"
	"lock-stock-v2/internal/domain/room/service"
	"lock-stock-v2/internal/domain/room_user/service"
	userRepository "lock-stock-v2/internal/domain/user/repository"
	"log"
	"net/http"
)

type RoomHandler struct {
	joinRoomService  *services.JoinRoomService
	roomUserService  *services.RoomUserService
	startGameService *service.StartGameService
	roomRepository   repository.RoomRepository
	userRepository   userRepository.UserRepository
}

func NewRoomHandler(
	u *services.JoinRoomService,
	roomRepository repository.RoomRepository,
	userRepository userRepository.UserRepository,
	roomUserService *services.RoomUserService,
	startGameService *service.StartGameService,
) *RoomHandler {
	return &RoomHandler{
		joinRoomService:  u,
		roomRepository:   roomRepository,
		roomUserService:  roomUserService,
		startGameService: startGameService,
		userRepository:   userRepository,
	}
}

func (h *RoomHandler) GetRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.roomRepository.GetPending()
	if err != nil {
		http.Error(w, "Failed to get rooms: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var responseData []RoomResponse
	for _, room := range rooms {
		var responseRoomUid = room.Uid()
		responseData = append(responseData, RoomResponse{
			RoomUid: &responseRoomUid,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		log.Printf("Ошибка при кодировании ответа: %v\n", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RoomHandler) StartGame(w http.ResponseWriter, r *http.Request, roomId string) {
	room, err := helpers.GetRoomById(h.roomRepository, roomId)
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

	req := service.StartGameRequest{
		Room: room,
		User: user,
	}
	if err := h.startGameService.StartGame(req); err != nil {
		http.Error(w, "Failed to start game: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Game started"})
}

func (h *RoomHandler) JoinRoom(w http.ResponseWriter, r *http.Request, roomId string, params JoinRoomParams) {
	room, err := helpers.GetRoomById(h.roomRepository, roomId)
	if err != nil {
		http.Error(w, err.Error(), err.(*helpers.RoomNotFoundError).Code)
		return
	}

	user, err := helpers.GetUserFromString(params.Authorization, h.userRepository)
	if err != nil {
		http.Error(w, err.Error(), err.(*helpers.UserNotFoundError).Code)
		return
	}
	domainReq := services.JoinRoomRequest{User: user, Room: room}
	if err := h.joinRoomService.JoinRoom(domainReq); err != nil {
		http.Error(w, "Failed to join room: "+err.Error(), http.StatusInternalServerError)
		return
	}

	roomUsers, err := h.roomUserService.GetUsersByRoom(room)
	if err != nil {
		http.Error(w, "Failed to get users in room", http.StatusInternalServerError)
		return
	}

	var response []JoinRoomResponse
	for _, ru := range roomUsers {
		roomFromRoomUser := ru.Room()
		u := ru.User()
		response = append(response, JoinRoomResponse{
			RoomId:      roomFromRoomUser.Uid(),
			UserId:      u.Uid(),
			UserName:    u.Name(),
			UserBalance: "5000",
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
