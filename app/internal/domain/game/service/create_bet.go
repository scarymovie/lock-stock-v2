package service

import (
	"lock-stock-v2/internal/domain/game/model"
	"lock-stock-v2/internal/domain/game/repository"
	"log"
)

type CreateBetService struct {
	betRepository repository.BetRepository
}

func NewCreateBetService(betRepository repository.BetRepository) *CreateBetService {
	return &CreateBetService{betRepository: betRepository}
}

func (cbs CreateBetService) CreateBet(player *model.Player, amount int, round *model.Round, position uint) (*model.Bet, error) {
	bet := model.NewBet(player, amount, round, position)
	err := cbs.betRepository.Save(bet)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return bet, nil
}
