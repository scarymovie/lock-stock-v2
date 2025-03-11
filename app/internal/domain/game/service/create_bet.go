package service

import (
	"encoding/json"
	"lock-stock-v2/internal/domain/game/model"
	"lock-stock-v2/internal/domain/game/repository"
	"lock-stock-v2/internal/websocket"
	"log"
)

type CreateBetService struct {
	betRepository repository.BetRepository
	webSocket     websocket.Manager
}

type NewBetMessage struct {
	Event string     `json:"event"`
	Body  NewBetBody `json:"body"`
}

type NewBetBody struct {
	UserID string `json:"userId"`
	Amount int    `json:"amount"`
}

func NewCreateBetService(betRepository repository.BetRepository, websocket websocket.Manager) *CreateBetService {
	return &CreateBetService{betRepository: betRepository, webSocket: websocket}
}

func (cbs CreateBetService) CreateBet(player *model.Player, amount int, round *model.Round) (*model.Bet, error) {
	bet := model.NewBet(player, amount, round)
	err := cbs.betRepository.Save(bet)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	message := NewBetMessage{
		Event: "new_bet",
		Body: NewBetBody{
			UserID: player.User().Uid(),
			Amount: amount,
		},
	}

	if err := cbs.sendWebSocketMessage(round.Game().Room().Uid(), message); err != nil {
		return nil, err
	}

	return bet, nil
}

func (cbs CreateBetService) sendWebSocketMessage(roomID string, message NewBetMessage) error {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal WebSocket message: %v\n", err)
		return err
	}

	log.Println(string(jsonMessage))
	cbs.webSocket.PublishToRoom(roomID, jsonMessage)
	return nil
}
