package repository

import (
	roomModel "lock-stock-v2/internal/domain/room/model"
	"lock-stock-v2/internal/domain/room_user/model"
)

type RoomUserRepository interface {
	Save(roomUser *model.RoomUser) error
	FindByRoom(room *roomModel.Room) ([]*model.RoomUser, error)
}
