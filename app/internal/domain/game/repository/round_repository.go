package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"lock-stock-v2/internal/domain/game/model"
)

type RoundRepository interface {
	FindByGame(game *model.LockStockGame) ([]*model.Round, error)
	FindLastByGame(game *model.LockStockGame) (*model.Round, error)
	Save(ctx context.Context, tx pgx.Tx, round *model.Round) error
}
