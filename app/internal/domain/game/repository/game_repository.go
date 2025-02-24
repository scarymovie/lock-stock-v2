package repository

import "lock-stock-v2/internal/domain/game/model"

type GameRepository interface {
	Save(game *model.LockStockGame)
}
