package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"lock-stock-v2/internal/domain/game/model"
	userModel "lock-stock-v2/internal/domain/user/model"
)

type GameRepository interface {
	Save(ctx context.Context, tx pgx.Tx, game *model.LockStockGame) error
	FindByUser(user *userModel.User) (*model.LockStockGame, error)
}
