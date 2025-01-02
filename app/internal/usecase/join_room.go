package usecase

import (
	"encoding/json"
	"lock-stock-v2/external/domain"
	"lock-stock-v2/external/usecase"
	internalDomain "lock-stock-v2/internal/domain"
	internalWs "lock-stock-v2/internal/websocket"
	"log"
)

type JoinRoomUsecase struct {
	roomUserRepository domain.RoomUserRepository
	webSocketManager   *internalWs.WebSocketManager // Зависимость для публикации
}

func NewJoinRoomUsecase(
	roomUserRepository domain.RoomUserRepository,
	webSocketManager *internalWs.WebSocketManager,
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

	log.Printf("usecase Player %s joined room %s\n", request.User.GetUserId(), request.Room.GetRoomUid())

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
