package service

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"lock-stock-v2/external/websocket"
	"lock-stock-v2/internal/domain/game/model"
	"lock-stock-v2/internal/domain/game/repository"
	"log"
)

type CreateRoundService struct {
	roundRepo            repository.RoundRepository
	createBetService     *CreateBetService
	createRoundPlayerLog *CreateRoundPlayerLog
	webSocket            websocket.Manager
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

func NewCreateRoundService(
	roundRepo repository.RoundRepository,
	createBetService *CreateBetService,
	webSocket websocket.Manager,
) *CreateRoundService {
	return &CreateRoundService{
		roundRepo:        roundRepo,
		createBetService: createBetService,
		webSocket:        webSocket,
	}
}

func (s *CreateRoundService) CreateRound(ctx context.Context, tx pgx.Tx, game *model.LockStockGame, players []*model.Player) error {
	rounds, _ := s.roundRepo.FindByGame(game)
	roundNumber := uint(1)
	if len(rounds) > 0 {
		roundNumber = uint(len(rounds) + 1)
	}

	roundId := "round-" + uuid.New().String()
	round := model.NewRound(roundId, &roundNumber, uint(roundCoefficient), 0, game)
	roundPrice := roundCoefficient * int(roundNumber)
	err := s.roundRepo.Save(ctx, tx, round)
	if err != nil {
		log.Printf("Error saving round: %s, %s", round.Uid(), err.Error())
		return err
	}

	var bets []*model.Bet
	for _, player := range players {
		newBalance := 0
		betValue := player.Balance()
		if roundPrice < player.Balance() {
			newBalance = player.Balance() - roundPrice
			betValue = roundPrice
		}
		bet, err := s.createBetService.CreateBet(ctx, tx, player, betValue, round)
		if err != nil {
			return err
		}
		bets = append(bets, bet)
		player.SetBalance(newBalance)
	}

	pot := 0
	for _, bet := range bets {
		pot += bet.Amount()
		round.SetPot(uint(pot))
	}

	err = s.roundRepo.Save(ctx, tx, round)
	if err != nil {
		log.Printf("Failed to save round player log: %v\n", err)
		return err
	}

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
