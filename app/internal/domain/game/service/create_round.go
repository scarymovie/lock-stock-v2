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

	round := model.NewRound(roundId, &roundNumber, uint(500), 0, game)
	roundPrice := roundCoefficient * int(roundNumber)
	s.roundRepo.Save(round)

	var bets []*model.Bet
	for i, player := range players {
		newBalance := 0
		betValue := player.Balance()
		if roundPrice < player.Balance() {
			newBalance = player.Balance() - roundPrice
			betValue = roundPrice
		}
		bet, _ := s.createBetService.CreateBet(player, betValue, round, uint(i+1))
		bets = append(bets, bet)
		player.SetBalance(newBalance)
	}

	pot := 0
	for _, bet := range bets {
		pot += bet.Amount()
		round.SetPot(uint(pot))
	}
	s.roundRepo.Save(round)

	body := map[string]interface{}{
		"roundNumber": roundNumber,
		"question":    NewQuestionMessage(round.Question()),
		"buyIn":       round.BuyIn(),
		"pot":         round.Pot(),
	}

	message := map[string]interface{}{
		"event": "round_started",
		"body":  body,
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
