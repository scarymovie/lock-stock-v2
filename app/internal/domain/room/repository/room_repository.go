package repository

import (
	"lock-stock-v2/internal/domain/room/model"
)

type RoomRepository interface {
	Save(room *model.Room) error
	UpdateRoomStatus(room *model.Room) error
	FindById(roomId string) (*model.Room, error)
	GetPending() ([]*model.Room, error)
}
