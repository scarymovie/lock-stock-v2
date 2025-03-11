package service

import (
	"encoding/json"
	"github.com/google/uuid"
	"lock-stock-v2/internal/domain/game/model"
	"lock-stock-v2/internal/domain/game/repository"
	"lock-stock-v2/internal/websocket"
	"log"
)

type CreateRoundService struct {
	roundRepo        repository.RoundRepository
	createBetService *CreateBetService
	webSocket        websocket.Manager
}

type QuestionMessage struct {
	Text  string   `json:"text"`
	Hints []string `json:"hints"`
}

type RoundStartedMessage struct {
	Event string                  `json:"event"`
	Body  RoundStartedMessageBody `json:"body"`
}

type RoundStartedMessageBody struct {
	RoundNumber uint            `json:"roundNumber"`
	Question    QuestionMessage `json:"question"`
	BuyIn       uint            `json:"buyIn"`
	Pot         uint            `json:"pot"`
}

const roundCoefficient = 500

func NewCreateRoundService(roundRepo repository.RoundRepository, createBetService *CreateBetService, webSocket websocket.Manager) *CreateRoundService {
	return &CreateRoundService{roundRepo: roundRepo, createBetService: createBetService, webSocket: webSocket}
}

func (s *CreateRoundService) CreateRound(game *model.LockStockGame, players []*model.Player) error {
	rounds, _ := s.roundRepo.FindByGame(game)
	roundNumber := uint(1)
	if len(rounds) > 0 {
		roundNumber = uint(len(rounds) + 1)
	}

	roundId := "round-" + uuid.New().String()
	round := model.NewRound(roundId, &roundNumber, uint(roundCoefficient), 0, game)
	roundPrice := roundCoefficient * int(roundNumber)
	s.roundRepo.Save(round)

	var bets []*model.Bet
	var roundPlayerLogs []*model.RoundPlayerLog
	for i, player := range players {
		newBalance := 0
		betValue := player.Balance()
		if roundPrice < player.Balance() {
			newBalance = player.Balance() - roundPrice
			betValue = roundPrice
		}
		bet, _ := s.createBetService.CreateBet(player, betValue, round)
		bets = append(bets, bet)

		roundPlayerLog := model.NewRoundPlayerLog(player, round, uint(betValue), uint(i)+1)
		roundPlayerLogs = append(roundPlayerLogs, roundPlayerLog)

		player.SetBalance(newBalance)
	}

	pot := 0
	for _, bet := range bets {
		pot += bet.Amount()
		round.SetPot(uint(pot))
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
		player := roundPlayerLog.Player()

		if betsValue < minBetsValue {
			minBetsValue = betsValue
			minNumber = number
			selectedPlayer = player
		} else if betsValue == minBetsValue && number < minNumber {
			minNumber = number
			selectedPlayer = player
		}

		if betsValue > maxBetsValue {
			maxBetsValue = betsValue
		}
	}

	if selectedPlayer != nil {
		round.SetPlayerTurn(selectedPlayer)
		round.SetMaxBet(maxBetsValue)
	}

	s.roundRepo.Save(round)

	return s.sendRoundStartedMessage(game, round)
}

func (s *CreateRoundService) sendRoundStartedMessage(game *model.LockStockGame, round *model.Round) error {
	message := RoundStartedMessage{
		Event: "round_started",
		Body: RoundStartedMessageBody{
			RoundNumber: *round.Number(),
			Question:    NewQuestionMessage(round.Question()),
			BuyIn:       round.BuyIn(),
			Pot:         round.Pot(),
		},
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal WebSocket message: %v\n", err)
		return err
	}

	log.Println(string(jsonMessage))
	s.webSocket.PublishToRoom(game.Room().Uid(), jsonMessage)
	return nil
}

func NewQuestionMessage(q *model.Question) QuestionMessage {
	hints := make([]string, len(q.Hints()))
	for i, hint := range q.Hints() {
		hints[i] = hint.Text()
	}
	return QuestionMessage{
		Text:  q.Text(),
		Hints: hints,
	}
}
