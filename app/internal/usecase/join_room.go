package usecase

import (
	"encoding/json"
	"lock-stock-v2/external/domain"
	"lock-stock-v2/external/usecase"
	externalWs "lock-stock-v2/external/websocket"
	internalDomain "lock-stock-v2/internal/domain"
	"log"
)

type JoinRoomUsecase struct {
	roomUserRepository domain.RoomUserRepository
	webSocketManager   externalWs.Manager
}

func NewJoinRoomUsecase(
	roomUserRepository domain.RoomUserRepository,
	webSocketManager externalWs.Manager,
) *JoinRoomUsecase {
	return &JoinRoomUsecase{
		roomUserRepository: roomUserRepository,
		webSocketManager:   webSocketManager,
	}
}

func (s JoinRoomUsecase) JoinRoom(request usecase.JoinRoomRequest) error {
	roomUser := internalDomain.RoomUser{}
	roomUser.SetRoom(request.Room)
	roomUser.SetUser(request.User)

	err := s.roomUserRepository.Save(&roomUser)
	if err != nil {
		return err
	}

	log.Printf("usecase Player %s joined room %s\n", request.User.GetUserUid(), request.Room.GetRoomUid())

	message := map[string]string{
		"event":   "user_joined",
		"user_id": roomUser.GetUser().GetUserUid(),
		"room_id": roomUser.GetRoom().GetRoomUid(),
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal WebSocket message: %v\n", err)
		return nil
	}

	log.Println(string(jsonMessage))
	s.webSocketManager.PublishToRoom(request.Room.GetRoomUid(), jsonMessage)

	return nil
}
