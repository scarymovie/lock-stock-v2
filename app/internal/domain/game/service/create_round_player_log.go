package service

import (
	"lock-stock-v2/internal/domain/game/model"
	"lock-stock-v2/internal/domain/game/repository"
)

type CreateRoundPlayerLog struct {
	roundPlayerLogRepository repository.RoundPlayerLogRepository
}

func NewCreateRoundPlayerLog(roundPlayerLogRepository repository.RoundPlayerLogRepository) *CreateRoundPlayerLog {
	return &CreateRoundPlayerLog{roundPlayerLogRepository: roundPlayerLogRepository}
}

func (l *CreateRoundPlayerLog) CreateRoundPlayerLog(player *model.Player, round *model.Round, amount uint, position uint) (*model.RoundPlayerLog, error) {
	roundPlayerLog := model.NewRoundPlayerLog(player, round, amount, position)
	err := l.roundPlayerLogRepository.Save(roundPlayerLog)
	if err != nil {
		return nil, err
	}
	return roundPlayerLog, nil
}
