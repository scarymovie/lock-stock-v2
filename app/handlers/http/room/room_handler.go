package room

import (
	"encoding/json"
	"errors"
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
		respondWithError(w, "Failed to get rooms", err, http.StatusInternalServerError)
		return
	}

	var responseData []RoomResponse
	for _, room := range rooms {
		var responseRoomUid = room.Uid()
		responseData = append(responseData, RoomResponse{
			RoomUid: &responseRoomUid,
		})
	}

	respondWithJSON(w, http.StatusOK, responseData)
}

func (h *RoomHandler) StartGame(w http.ResponseWriter, r *http.Request, roomId string) {
	room, err := helpers.GetRoomById(h.roomRepository, roomId)
	if err != nil {
		var roomErr *helpers.RoomNotFoundError
		ok := errors.As(err, &roomErr)
		if ok {
			respondWithError(w, err.Error(), nil, roomErr.Code)
		} else {
			respondWithError(w, "Error getting room", err, http.StatusInternalServerError)
		}
		return
	}

	user, err := helpers.GetUserFromRequest(r)
	if err != nil {
		var userErr *helpers.UserNotFoundError
		ok := errors.As(err, &userErr)
		if ok {
			respondWithError(w, err.Error(), nil, userErr.Code)
		} else {
			respondWithError(w, "Error getting user", err, http.StatusInternalServerError)
		}
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
		respondWithError(w, "Failed to start game", err, http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Game started"})
}

func (h *RoomHandler) JoinRoom(w http.ResponseWriter, r *http.Request, roomId string, params JoinRoomParams) {
	room, err := helpers.GetRoomById(h.roomRepository, roomId)
	if err != nil {
		var roomErr *helpers.RoomNotFoundError
		ok := errors.As(err, &roomErr)
		if ok {
			respondWithError(w, err.Error(), nil, roomErr.Code)
		} else {
			respondWithError(w, "Error getting room", err, http.StatusInternalServerError)
		}
		return
	}

	user, err := helpers.GetUserFromString(params.Authorization, h.userRepository)
	if err != nil {
		userErr, ok := err.(*helpers.UserNotFoundError)
		if ok {
			respondWithError(w, err.Error(), nil, userErr.Code)
		} else {
			respondWithError(w, "Error getting user", err, http.StatusInternalServerError)
		}
		return
	}

	domainReq := services.JoinRoomRequest{User: user, Room: room}
	if err := h.joinRoomService.JoinRoom(domainReq); err != nil {
		respondWithError(w, "Failed to join room", err, http.StatusInternalServerError)
		return
	}

	roomUsers, err := h.roomUserService.GetUsersByRoom(room)
	if err != nil {
		respondWithError(w, "Failed to get users in room", err, http.StatusInternalServerError)
		return
	}

	response := make([]JoinRoomResponse, 0, len(roomUsers))
	for _, ru := range roomUsers {
		response = append(response, JoinRoomResponse{
			RoomId:   ru.Room().Uid(),
			UserId:   ru.User().Uid(),
			UserName: ru.User().Name(),
		})
	}

	respondWithJSON(w, http.StatusOK, response)
}

func respondWithError(w http.ResponseWriter, message string, err error, statusCode int) {
	errorMessage := message
	if err != nil {
		errorMessage += ": " + err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": errorMessage})
}

func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Ошибка при кодировании ответа: %v\n", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
