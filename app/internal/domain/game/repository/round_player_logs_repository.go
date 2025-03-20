package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"lock-stock-v2/internal/domain/game/model"
	userModel "lock-stock-v2/internal/domain/user/model"
)

type RoundPlayerLogRepository interface {
	FindByRound(round *model.Round) ([]*model.RoundPlayerLog, error)
	FindByRoundAndUser(round *model.Round, user *userModel.User) (*model.RoundPlayerLog, error)
	Save(ctx context.Context, tx pgx.Tx, log *model.RoundPlayerLog) error
}
