package service

import (
	"context"
	"github.com/jackc/pgx/v5"
	"lock-stock-v2/internal/domain/game/model"
	"lock-stock-v2/internal/domain/game/repository"
)

type CreateRoundPlayerLog struct {
	roundPlayerLogRepository repository.RoundPlayerLogRepository
}

func NewCreateRoundPlayerLog(roundPlayerLogRepository repository.RoundPlayerLogRepository) *CreateRoundPlayerLog {
	return &CreateRoundPlayerLog{roundPlayerLogRepository: roundPlayerLogRepository}
}

func (l *CreateRoundPlayerLog) CreateRoundPlayerLog(ctx context.Context, tx pgx.Tx, player *model.Player, round *model.Round, amount uint, position uint) (*model.RoundPlayerLog, error) {
	roundPlayerLog := model.NewRoundPlayerLog(player, round, position, amount, nil)
	err := l.roundPlayerLogRepository.Save(ctx, tx, roundPlayerLog)
	if err != nil {
		return nil, err
	}
	return roundPlayerLog, nil
}
