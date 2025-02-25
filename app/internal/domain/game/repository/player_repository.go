package repository

import "lock-stock-v2/internal/domain/game/model"

type PlayerRepository interface {
	Save(player *model.Player) error
}
