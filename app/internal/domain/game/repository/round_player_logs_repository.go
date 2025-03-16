package repository

import (
	"lock-stock-v2/internal/domain/game/model"
	userModel "lock-stock-v2/internal/domain/user/model"
)

type RoundPlayerLogRepository interface {
	FindByRound(round *model.Round) ([]*model.RoundPlayerLog, error)
	FindByRoundAndUser(round *model.Round, user *userModel.User) (*model.RoundPlayerLog, error)
	Save(log *model.RoundPlayerLog) error
}
