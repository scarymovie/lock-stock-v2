package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"lock-stock-v2/internal/domain/game/model"
	userModel "lock-stock-v2/internal/domain/user/model"
)

type RoundPlayerLogRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRoundPlayerLogRepository(db *pgxpool.Pool) *RoundPlayerLogRepository {
	return &RoundPlayerLogRepository{db: db}
}

func (repo *RoundPlayerLogRepository) FindByRound(round *model.Round) ([]*model.RoundPlayerLog, error) {
	return []*model.RoundPlayerLog{}, nil
}
func (repo *RoundPlayerLogRepository) FindByRoundAndUser(round *model.Round, user *userModel.User) (*model.RoundPlayerLog, error) {
	return nil, nil
}
func (repo *RoundPlayerLogRepository) Save(log *model.RoundPlayerLog) error {
	return nil
}
