package inMemory

import (
	"errors"
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

func (repo *RoomRepository) FindById(roomId string) (*model.Room, error) {
	room, exists := repo.rooms[roomId]
	if !exists {
		return nil, errors.New("room not found")
	}
	return room, nil
}

func (repo *RoomRepository) Save(room *model.Room) error {
	repo.rooms[room.Uid()] = room
	return nil
}

func (repo *RoomRepository) Count() int {
	return len(repo.rooms)
}

func (repo *RoomRepository) UpdateRoomStatus(room *model.Room) error {
	return nil
}

func (repo *RoomRepository) GetPending() ([]*model.Room, error) {
	return nil, nil
}
