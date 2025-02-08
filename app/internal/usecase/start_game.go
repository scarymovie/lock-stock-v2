package usecase

import (
	"encoding/json"
	"lock-stock-v2/external/domain"
	"lock-stock-v2/external/usecase"
	"lock-stock-v2/external/websocket"
	"log"
)

type StartGameUsecase struct {
	roomRepository domain.RoomRepository
	webSocket      websocket.Manager
}

func NewStartGameUsecase(
	roomRepository domain.RoomRepository,
	webSocket websocket.Manager,
) *StartGameUsecase {
	return &StartGameUsecase{roomRepository: roomRepository, webSocket: webSocket}
}

func (uc *StartGameUsecase) StartGame(req usecase.StartGameRequest) error {

	req.Room.SetStatus(domain.StatusStarted)
	if err := uc.roomRepository.UpdateRoomStatus(req.Room); err != nil {
		log.Println("Failed to update room status:", err)
		return err
	}

	message := map[string]string{
		"event":   "game_started",
		"roomUid": req.Room.Uid(),
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
