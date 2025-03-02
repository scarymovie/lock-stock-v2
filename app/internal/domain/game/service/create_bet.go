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

func NewCreateBetService(betRepository repository.BetRepository, websocket websocket.Manager) *CreateBetService {
	return &CreateBetService{betRepository: betRepository, webSocket: websocket}
}

func (cbs CreateBetService) CreateBet(player *model.Player, amount int, round *model.Round, position uint) (*model.Bet, error) {
	bet := model.NewBet(player, amount, round, position)
	err := cbs.betRepository.Save(bet)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	body := map[string]interface{}{
		"playerId": player.Uid(),
		"amount":   amount,
	}
	message := map[string]interface{}{
		"event": "new_bet",
		"body":  body,
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal WebSocket message: %v\n", err)
		return nil, err
	}

	log.Println(string(jsonMessage))
	cbs.webSocket.PublishToRoom(round.Game().Room().Uid(), jsonMessage)

	return bet, nil
}
