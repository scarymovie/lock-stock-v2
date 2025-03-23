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

type CreateBetService struct {
	betRepository            repository.BetRepository
	webSocket                websocket.Manager
	roundPlayerLogRepository repository.RoundPlayerLogRepository
	createRoundPlayerLog     *CreateRoundPlayerLog
	roundObserver            *RoundObserver
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
var ErrRoundPlayerLogsByRoundNotFound = errors.New("round player logs by round not found")

func NewCreateBetService(
	betRepository repository.BetRepository,
	websocket websocket.Manager,
	roundPlayerLogRepository repository.RoundPlayerLogRepository,
	createRoundPlayerLog *CreateRoundPlayerLog,
	roundObserver *RoundObserver,
) *CreateBetService {
	return &CreateBetService{
		betRepository:            betRepository,
		webSocket:                websocket,
		roundPlayerLogRepository: roundPlayerLogRepository,
		createRoundPlayerLog:     createRoundPlayerLog,
		roundObserver:            roundObserver,
	}
}

func (cbs *CreateBetService) CreateBet(ctx context.Context, tx pgx.Tx, player *model.Player, amount int, round *model.Round) (*model.Bet, error) {
	if player == nil {
		log.Println("Player is nil")
		return nil, ErrPlayerNotFound
	}
	bet := model.NewBet(player, amount, round)
	err := cbs.betRepository.Save(ctx, tx, bet)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	roundPlayerLogs, err := cbs.roundPlayerLogRepository.FindByRound(round)
	if err != nil {
		log.Println(err)
		return nil, ErrRoundPlayerLogsByRoundNotFound
	}

	if len(roundPlayerLogs) == 0 {
		_, err = cbs.createRoundPlayerLog.CreateRoundPlayerLog(ctx, tx, player, round, uint(amount), uint(0)+1)

		if err != nil {
			log.Printf("Failed to create round player log: %v\n", err)
			return nil, err
		}
	}

	for _, roundPlayerLog := range roundPlayerLogs {
		if roundPlayerLog.Player().User().Uid() == player.User().Uid() {
			roundPlayerLog.SetBetsValue(roundPlayerLog.BetsValue() + uint(amount))
			err = cbs.roundPlayerLogRepository.Save(ctx, tx, roundPlayerLog)
			if err != nil {
				log.Println("Failed to save round player log", err.Error())
				return nil, err
			}
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

	if err = cbs.sendWebSocketMessage(round.Game().Room().Uid(), message); err != nil {
		return nil, err
	}

	err = cbs.roundObserver.ObserveRoundState(round)
	if nil != err {
		log.Println("Observer error")
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
