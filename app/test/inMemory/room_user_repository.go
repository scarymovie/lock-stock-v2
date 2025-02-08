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

	key := fmt.Sprintf("%s:%s", ru.GetRoom().Uid(), ru.GetUser().Uid())

	repo.roomUsers[key] = ru
	return nil
}

func (repo *RoomUserRepository) FindByRoom(room api.Room) ([]api.RoomUser, error) {
	roomUid := room.Uid()

	var result []api.RoomUser
	for _, ru := range repo.roomUsers {
		if ru.GetRoom().Uid() == roomUid {
			result = append(result, ru)
		}
	}
	if len(result) == 0 {
		return nil, errors.New("no users found in room")
	}

	return result, nil
}

func (repo *RoomUserRepository) Count() int {
	return len(repo.roomUsers)
}
