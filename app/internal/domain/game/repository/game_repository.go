package repository

import (
	"lock-stock-v2/internal/domain/game/model"
	userModel "lock-stock-v2/internal/domain/user/model"
)

type GameRepository interface {
	Save(game *model.LockStockGame) error
	FindByUser(user *userModel.User) (*model.LockStockGame, error)
}
