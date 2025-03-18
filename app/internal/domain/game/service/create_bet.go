package service

import (
	"encoding/json"
	"errors"
	"lock-stock-v2/internal/domain/game/model"
	"lock-stock-v2/internal/domain/game/repository"
	"lock-stock-v2/internal/websocket"
	"log"
)

type CreateBetService struct {
	betRepository            repository.BetRepository
	webSocket                websocket.Manager
	roundPlayerLogRepository repository.RoundPlayerLogRepository
}

type NewBetMessage struct {
	Event string     `json:"event"`
	Body  NewBetBody `json:"body"`
}

type NewBetBody struct {
	UserID           string `json:"userId"`
	Amount           int    `json:"amount"`
	NextPlayerTurnID string `json:"nextPlayerTurnID"`
	MaxBet           uint   `json:"maxBet"`
}

var ErrPlayerNotFound = errors.New("player not found")

func NewCreateBetService(betRepository repository.BetRepository, websocket websocket.Manager, roundPlayerLogRepository repository.RoundPlayerLogRepository) *CreateBetService {
	return &CreateBetService{betRepository: betRepository, webSocket: websocket, roundPlayerLogRepository: roundPlayerLogRepository}
}

func (cbs *CreateBetService) CreateBet(player *model.Player, amount int, round *model.Round) (*model.Bet, error) {
	if player == nil {
		log.Println("Player is nil")
		return nil, ErrPlayerNotFound
	}
	bet := model.NewBet(player, amount, round)
	err := cbs.betRepository.Save(bet)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	roundPlayerLogs, _ := cbs.roundPlayerLogRepository.FindByRound(round)

	for _, roundPlayerLog := range roundPlayerLogs {
		if roundPlayerLog.Player().User().Uid() == player.User().Uid() {
			roundPlayerLog.SetBetsValue(roundPlayerLog.BetsValue() + uint(amount))
		}
	}

	var (
		minBetsValue   = ^uint(0)
		maxBetsValue   uint
		minNumber      = ^uint(0)
		selectedPlayer *model.Player
	)

	for _, roundPlayerLog := range roundPlayerLogs {
		betsValue := roundPlayerLog.BetsValue()
		number := roundPlayerLog.Number()
		roundPlayerLogPlayer := roundPlayerLog.Player()

		if betsValue < minBetsValue {
			minBetsValue = betsValue
			minNumber = number
			selectedPlayer = roundPlayerLogPlayer
		} else if betsValue == minBetsValue && number < minNumber {
			minNumber = number
			selectedPlayer = roundPlayerLogPlayer
		}

		if betsValue > maxBetsValue {
			maxBetsValue = betsValue
		}
	}

	if selectedPlayer != nil {
		round.SetPlayerTurn(selectedPlayer)
		round.SetMaxBet(maxBetsValue)
	}

	message := NewBetMessage{
		Event: "new_bet",
		Body: NewBetBody{
			UserID: player.User().Uid(),
			Amount: amount,
			MaxBet: round.MaxBet(),
		},
	}

	if err := cbs.sendWebSocketMessage(round.Game().Room().Uid(), message); err != nil {
		return nil, err
	}

	return bet, nil
}

func (cbs *CreateBetService) sendWebSocketMessage(roomID string, message NewBetMessage) error {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal WebSocket message: %v\n", err)
		return err
	}

	log.Println(string(jsonMessage))
	cbs.webSocket.PublishToRoom(roomID, jsonMessage)
	return nil
}
