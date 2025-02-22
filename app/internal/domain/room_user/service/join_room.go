package services

import (
	"encoding/json"
	roomModel "lock-stock-v2/internal/domain/room/model"
	roomUserModel "lock-stock-v2/internal/domain/room_user/model"
	"lock-stock-v2/internal/domain/room_user/repository"
	userModel "lock-stock-v2/internal/domain/user/model"
	internalWs "lock-stock-v2/internal/websocket"
	"log"
)

type JoinRoomService struct {
	roomUserRepository repository.RoomUserRepository
	webSocketManager   internalWs.Manager
}

type JoinRoomRequest struct {
	User *userModel.User
	Room *roomModel.Room
}

func NewJoinRoomService(
	roomUserRepository repository.RoomUserRepository,
	webSocketManager internalWs.Manager,
) *JoinRoomService {
	return &JoinRoomService{
		roomUserRepository: roomUserRepository,
		webSocketManager:   webSocketManager,
	}
}

func (s JoinRoomService) JoinRoom(request JoinRoomRequest) error {
	roomUser := roomUserModel.NewRoomUser(request.Room, request.User)

	err := s.roomUserRepository.Save(roomUser)
	if err != nil {
		return err
	}

	log.Printf("usecase Player %s joined room %s\n", request.User.Uid(), request.Room.Uid())

	message := map[string]string{
		"event":    "user_joined",
		"userId":   roomUser.User().Uid(),
		"userName": roomUser.User().Name(),
		"roomId":   roomUser.Room().Uid(),
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal WebSocket message: %v\n", err)
		return nil
	}

	log.Println(string(jsonMessage))
	s.webSocketManager.PublishToRoom(request.Room.Uid(), jsonMessage)

	return nil
}
