package service

import (
	"encoding/json"
	"log"

	"lock-stock-v2/external/websocket"
	"lock-stock-v2/internal/domain/game/model"
	"lock-stock-v2/internal/domain/game/repository"
)

type RoundObserver struct {
	roundPlayerLogRepo repository.RoundPlayerLogRepository
	webSocket          websocket.Manager
}

type AllAnswersReceivedMessage struct {
	Event string                 `json:"event"`
	Body  AllAnswersReceivedBody `json:"body"`
}

type AllAnswersReceivedBody struct {
	RoundUID string `json:"roundUid"`
}

type BettingFinishedMessage struct {
	Event string              `json:"event"`
	Body  BettingFinishedBody `json:"body"`
}

type BettingFinishedBody struct {
	RoundUID string `json:"roundUid"`
	MinBet   uint   `json:"minBet"`
}

func NewRoundObserver(repo repository.RoundPlayerLogRepository, ws websocket.Manager) *RoundObserver {
	return &RoundObserver{
		roundPlayerLogRepo: repo,
		webSocket:          ws,
	}
}

func (ro *RoundObserver) ObserveRoundState(round *model.Round) {
	roundPlayerLogs, err := ro.roundPlayerLogRepo.FindByRound(round)
	if err != nil {
		log.Printf("Observer: Failed to get RoundPlayerLogs: %v", err)
		return
	}

	allAnswered := true
	for _, roundPlayerLog := range roundPlayerLogs {
		if roundPlayerLog.Answer() == nil {
			allAnswered = false
			break
		}
	}

	if allAnswered {
		ro.sendWebSocketMessage(round.Game().Room().Uid(), AllAnswersReceivedMessage{
			Event: "all_answers_received",
			Body: AllAnswersReceivedBody{
				RoundUID: round.Uid(),
			},
		})
	}

	allEqual := true
	var targetBet uint
	first := true

	for _, roundPlayerLog := range roundPlayerLogs {
		if roundPlayerLog.Player().Status() == model.StatusLost {
			continue
		}
		if first {
			targetBet = roundPlayerLog.BetsValue()
			first = false
		} else if roundPlayerLog.BetsValue() != targetBet {
			allEqual = false
			break
		}
	}

	if allEqual && !first {
		ro.sendWebSocketMessage(round.Game().Room().Uid(), BettingFinishedMessage{
			Event: "betting_finished",
			Body: BettingFinishedBody{
				RoundUID: round.Uid(),
				MinBet:   targetBet,
			},
		})
	}
}

func (ro *RoundObserver) sendWebSocketMessage(roomID string, data any) {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Observer: Failed to marshal WS message: %v", err)
		return
	}
	ro.webSocket.PublishToRoom(roomID, bytes)
}
