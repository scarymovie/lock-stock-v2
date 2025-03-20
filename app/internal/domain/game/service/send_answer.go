package service

import (
	"encoding/json"
	"errors"
	"lock-stock-v2/external/websocket"
	"lock-stock-v2/internal/domain/game/model"
	"lock-stock-v2/internal/domain/game/repository"
	"log"
)

type SendAnswer struct {
	roundPlayerLogRepository repository.RoundPlayerLogRepository
	webSocket                websocket.Manager
}

func NewSendAnswer(
	roundPlayerLogRepository repository.RoundPlayerLogRepository,
	webSocket websocket.Manager,
) *SendAnswer {
	return &SendAnswer{
		roundPlayerLogRepository: roundPlayerLogRepository,
		webSocket:                webSocket,
	}
}

type NewAnswerMessage struct {
	Event string        `json:"event"`
	Body  NewAnswerBody `json:"body"`
}

type NewAnswerBody struct {
	UserId string `json:"userId"`
}

var RoundPlayerLogIsNil = errors.New("roundPlayerLog is nil")

func (s *SendAnswer) SendAnswer(roundPlayerLog *model.RoundPlayerLog, answer uint) error {
	if nil == roundPlayerLog {
		log.Println("roundPlayerLog is nil")
		return RoundPlayerLogIsNil
	}
	roundPlayerLog.SetAnswer(&answer)
	err := s.roundPlayerLogRepository.Save(roundPlayerLog)
	if err != nil {
		return err
	}
	message := NewAnswerMessage{
		Event: "new_answer",
		Body: NewAnswerBody{
			UserId: roundPlayerLog.Player().User().Uid(),
		},
	}

	if err = s.sendAnswerWebSocketMessage(roundPlayerLog.Round().Game().Room().Uid(), message); err != nil {
		return err
	}
	return nil
}

func (s *SendAnswer) sendAnswerWebSocketMessage(roomID string, message NewAnswerMessage) error {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal WebSocket message: %v\n", err)
		return err
	}

	log.Println(string(jsonMessage))
	s.webSocket.PublishToRoom(roomID, jsonMessage)
	return nil
}
