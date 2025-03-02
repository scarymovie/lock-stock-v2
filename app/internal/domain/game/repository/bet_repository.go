package repository

import "lock-stock-v2/internal/domain/game/model"

type BetRepository interface {
	Save(bet *model.Bet) error
	FindByRound(round *model.Round) ([]*model.Bet, error)
}
