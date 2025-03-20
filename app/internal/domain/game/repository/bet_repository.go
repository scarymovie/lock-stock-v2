package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"lock-stock-v2/internal/domain/game/model"
)

type BetRepository interface {
	Save(ctx context.Context, tx pgx.Tx, bet *model.Bet) error
	FindByRound(round *model.Round) ([]*model.Bet, error)
}
