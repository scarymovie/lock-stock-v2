package repository

import (
	"errors"
	"fmt"
	api "lock-stock-v2/external/domain"
	"lock-stock-v2/internal/domain"
)

type InMemoryRoomUserRepository struct {
	roomUsers map[string]*domain.RoomUser
}

func NewInMemoryRoomUserRepository() *InMemoryRoomUserRepository {
	return &InMemoryRoomUserRepository{
		roomUsers: make(map[string]*domain.RoomUser),
	}
}

func (repo *InMemoryRoomUserRepository) Save(roomUser api.RoomUser) error {
	ru, ok := roomUser.(*domain.RoomUser)
	if !ok {
		return errors.New("invalid RoomUser type")
	}

	// Генерация ключа по roomID и userID
	key := fmt.Sprintf("%s:%s", ru.GetRoom().GetRoomId(), ru.GetUser().GetUserId())

	repo.roomUsers[key] = ru
	return nil
}
