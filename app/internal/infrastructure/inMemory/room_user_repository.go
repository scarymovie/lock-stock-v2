package inMemory

import (
	"errors"
	"fmt"
	api "lock-stock-v2/external/domain"
	"lock-stock-v2/internal/domain"
)

type RoomUserRepository struct {
	roomUsers map[string]*domain.RoomUser
}

func NewInMemoryRoomUserRepository() *RoomUserRepository {
	return &RoomUserRepository{
		roomUsers: make(map[string]*domain.RoomUser),
	}
}

func (repo *RoomUserRepository) Save(roomUser api.RoomUser) error {
	ru, ok := roomUser.(*domain.RoomUser)
	if !ok {
		return errors.New("invalid RoomUser type")
	}

	// Генерация ключа по roomID и userID
	key := fmt.Sprintf("%s:%s", ru.GetRoom().GetRoomId(), ru.GetUser().GetUserId())

	repo.roomUsers[key] = ru
	return nil
}
