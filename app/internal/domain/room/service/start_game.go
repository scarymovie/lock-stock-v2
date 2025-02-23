package service

import (
	"encoding/json"
	roomModel "lock-stock-v2/internal/domain/room/model"
	"lock-stock-v2/internal/domain/room/repository"
	roomUserRepository "lock-stock-v2/internal/domain/room_user/repository"
	userModel "lock-stock-v2/internal/domain/user/model"
	"lock-stock-v2/internal/websocket"
	"log"
)

type StartGameService struct {
	roomRepository     repository.RoomRepository
	webSocket          websocket.Manager
	roomUserRepository roomUserRepository.RoomUserRepository
}

type StartGameRequest struct {
	Room *roomModel.Room
	User *userModel.User
}

func NewStartGameService(
	roomRepository repository.RoomRepository,
	roomUserRepository roomUserRepository.RoomUserRepository,
	webSocket websocket.Manager,
) *StartGameService {
	return &StartGameService{roomRepository: roomRepository, webSocket: webSocket, roomUserRepository: roomUserRepository}
}

func (uc *StartGameService) StartGame(req StartGameRequest) error {

	req.Room.SetStatus(roomModel.StatusStarted)
	if err := uc.roomRepository.UpdateRoomStatus(req.Room); err != nil {
		log.Println("Failed to update room status:", err)
		return err
	}

	message := map[string]interface{}{
		"event":            "game_started",
		"roomUid":          req.Room.Uid(),
		"questionDuration": "60",
		"actionDuration":   "30",
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal WebSocket message: %v\n", err)
		return nil
	}

	log.Println(string(jsonMessage))

	uc.webSocket.PublishToRoom(req.Room.Uid(), jsonMessage)

	return nil
}
