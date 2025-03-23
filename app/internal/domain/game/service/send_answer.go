package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"lock-stock-v2/external/websocket"
	"lock-stock-v2/internal/domain/game/model"
	"lock-stock-v2/internal/domain/game/repository"
	"log"
)

type SendAnswer struct {
	roundPlayerLogRepository repository.RoundPlayerLogRepository
	webSocket                websocket.Manager
	roundObserver            *RoundObserver
}

func NewSendAnswer(
	roundPlayerLogRepository repository.RoundPlayerLogRepository,
	webSocket websocket.Manager,
	roundObserver *RoundObserver,
) *SendAnswer {
	return &SendAnswer{
		roundPlayerLogRepository: roundPlayerLogRepository,
		webSocket:                webSocket,
		roundObserver:            roundObserver,
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

func (s *SendAnswer) SendAnswer(ctx context.Context, tx pgx.Tx, roundPlayerLog *model.RoundPlayerLog, answer uint) error {
	if nil == roundPlayerLog {
		log.Println("roundPlayerLog is nil")
		return RoundPlayerLogIsNil
	}
	roundPlayerLog.SetAnswer(&answer)
	err := s.roundPlayerLogRepository.Save(ctx, tx, roundPlayerLog)
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
	s.roundObserver.ObserveRoundState(roundPlayerLog.Round())

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
