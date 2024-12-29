package repository

import (
	"errors"
	api "lock-stock-v2/external/domain"
	"lock-stock-v2/internal/domain"
)

type InMemoryRoomRepository struct {
	rooms map[string]*domain.Room
}

func NewInMemoryRoomRepository() *InMemoryRoomRepository {
	return &InMemoryRoomRepository{
		rooms: map[string]*domain.Room{
			"room1": {Id: "room1"},
			"room2": {Id: "room2"},
			"room3": {Id: "room3"},
		},
	}
}

func (repo *InMemoryRoomRepository) FindById(roomId string) (api.Room, error) {
	room, exists := repo.rooms[roomId]
	if !exists {
		return nil, errors.New("room not found")
	}
	return room, nil
}

func (repo *InMemoryRoomRepository) Save(room api.Room) error {
	r, ok := room.(*domain.Room)
	if !ok {
		return errors.New("invalid room type")
	}

	repo.rooms[r.GetId()] = r
	return nil
}
