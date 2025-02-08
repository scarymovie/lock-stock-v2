package inMemory

import (
	"errors"
	api "lock-stock-v2/external/domain"
	"lock-stock-v2/internal/domain/room/model"
)

type RoomRepository struct {
	rooms map[string]*model.Room
}

func NewInMemoryRoomRepository() *RoomRepository {
	return &RoomRepository{
		rooms: make(map[string]*model.Room),
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
	r, ok := room.(*model.Room)
	if !ok {
		return errors.New("invalid room type")
	}

	repo.rooms[r.Uid()] = r
	return nil
}

func (repo *RoomRepository) Count() int {
	return len(repo.rooms)
}
