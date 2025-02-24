package service

import (
	"lock-stock-v2/internal/domain/game/model"
	"lock-stock-v2/internal/domain/game/repository"
)

type CreateRoundService struct {
	roundRepo repository.RoundRepository
}

func NewCreateRoundService() *CreateRoundService {
	return &CreateRoundService{}
}

func (s *CreateRoundService) CreateRound(game *model.LockStockGame, players []*model.Player) {
	rounds, _ := s.roundRepo.FindByGame(game)
	roundNumber := uint(1)
	if len(rounds) > 0 {
		roundNumber = uint(len(rounds) + 1)
	}

	hints := []*model.Hint{
		model.NewHint("Водолей является таким по счёту знаком Зодиака."),
		model.NewHint("Именно этого числа в России отмечается день трезвости."),
	}
	question := model.NewQuestion("Сколько голов забила Сборная России на чемпионате мира по футболу 2018 года?", hints)

	round := model.NewRound(&roundNumber, question, uint(500), game)
	s.roundRepo.Save(round)
}
