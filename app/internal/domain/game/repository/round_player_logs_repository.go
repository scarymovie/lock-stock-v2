package repository

import "lock-stock-v2/internal/domain/game/model"

type RoundPlayerLogRepository interface {
	FindByRound(round *model.Round) ([]*model.RoundPlayerLog, error)
}
