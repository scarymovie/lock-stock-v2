package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"lock-stock-v2/internal/domain/room/model"
)

type RoomRepository interface {
	Save(room *model.Room) error
	UpdateRoomStatus(ctx context.Context, tx pgx.Tx, room *model.Room) error
	FindById(roomId string) (*model.Room, error)
	GetPending() ([]*model.Room, error)
}
