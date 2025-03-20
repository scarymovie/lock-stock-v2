package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"lock-stock-v2/internal/domain/game/model"
	roomModel "lock-stock-v2/internal/domain/room/model"
	userModel "lock-stock-v2/internal/domain/user/model"
)

type PlayerRepository interface {
	Save(ctx context.Context, tx pgx.Tx, player *model.Player) error
	FindByUserAndRoom(user *userModel.User, room *roomModel.Room) (*model.Player, error)
}
