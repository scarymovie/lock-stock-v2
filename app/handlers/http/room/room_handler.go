package room

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
	helpers2 "lock-stock-v2/handlers/http/helpers"
	"lock-stock-v2/internal/domain/room/repository"
	"lock-stock-v2/internal/domain/room/service"
	"lock-stock-v2/internal/domain/room_user/service"
	"log"
	"net/http"
)

type RoomHandler struct {
	joinRoomService  *services.JoinRoomService
	roomRepository   repository.RoomRepository
	roomUserService  *services.RoomUserService
	startGameService *service.StartGameService
}

func NewRoomHandler(
	u *services.JoinRoomService,
	roomRepository repository.RoomRepository,
	roomUserService *services.RoomUserService,
	startGameService *service.StartGameService,
) *RoomHandler {
	return &RoomHandler{
		joinRoomService:  u,
		roomRepository:   roomRepository,
		roomUserService:  roomUserService,
		startGameService: startGameService,
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
		var responseRoomUid = uuid.MustParse(room.Uid())
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

func (h *RoomHandler) StartGame(w http.ResponseWriter, r *http.Request, roomId openapitypes.UUID) {
	roomID := chi.URLParam(r, "roomId")
	room, err := helpers2.GetRoomById(h.roomRepository, roomID)
	if err != nil {
		http.Error(w, err.Error(), err.(*helpers2.RoomNotFoundError).Code)
		return
	}

	user, err := helpers2.GetUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), err.(*helpers2.UserNotFoundError).Code)
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

func (h *RoomHandler) PostRoomJoinRoomId(w http.ResponseWriter, r *http.Request, roomId string) {
	var req JoinRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	room, err := helpers2.GetRoomById(h.roomRepository, roomId)
	if err != nil {
		http.Error(w, err.Error(), err.(*helpers2.RoomNotFoundError).Code)
		return
	}

	user, err := helpers2.GetUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), err.(*helpers2.UserNotFoundError).Code)
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
			RoomUid:  ptr(roomFromRoomUser.Uid()),
			UserUid:  ptr(u.Uid()),
			UserName: ptr(u.Name()),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func ptr(s string) *string {
	return &s
}
