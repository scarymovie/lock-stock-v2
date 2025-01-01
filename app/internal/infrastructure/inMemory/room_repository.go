package inMemory

import (
	"errors"
	api "lock-stock-v2/external/domain"
	"lock-stock-v2/internal/domain"
)

type RoomRepository struct {
	rooms map[string]*domain.Room
}

func NewInMemoryRoomRepository() *RoomRepository {
	return &RoomRepository{
		rooms: map[string]*domain.Room{
			"room1": {Id: "room1"},
			"room2": {Id: "room2"},
			"room3": {Id: "room3"},
		},
	}
}

func (repo *RoomRepository) FindById(roomId string) (api.Room, error) {
	room, exists := repo.rooms[roomId]
	if !exists {
		return nil, errors.New("room not found")
	}
	return room, nil
}

func (repo *RoomRepository) Save(room api.Room) error {
	r, ok := room.(*domain.Room)
	if !ok {
		return errors.New("invalid room type")
	}

	repo.rooms[r.GetRoomUid()] = r
	return nil
}
