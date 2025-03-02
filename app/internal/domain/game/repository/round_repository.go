package repository

import "lock-stock-v2/internal/domain/game/model"

type RoundRepository interface {
	FindByGame(game *model.LockStockGame) ([]*model.Round, error)
	FindLastByGame(game *model.LockStockGame) (*model.Round, error)
	Save(round *model.Round) error
}
