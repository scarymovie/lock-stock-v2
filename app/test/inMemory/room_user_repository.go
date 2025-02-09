package inMemory

import (
	"errors"
	"fmt"
	roomModel "lock-stock-v2/internal/domain/room/model"
	"lock-stock-v2/internal/domain/room_user/model"
)

type RoomUserRepository struct {
	roomUsers map[string]*model.RoomUser
}

func NewInMemoryRoomUserRepository() *RoomUserRepository {
	return &RoomUserRepository{
		roomUsers: make(map[string]*model.RoomUser),
	}
}

func (repo *RoomUserRepository) Save(roomUser *model.RoomUser) error {
	key := fmt.Sprintf("%s:%s", roomUser.Room().Uid(), roomUser.User().Uid())

	repo.roomUsers[key] = roomUser
	return nil
}

func (repo *RoomUserRepository) FindByRoom(room *roomModel.Room) ([]*model.RoomUser, error) {
	roomUid := room.Uid()

	var result []*model.RoomUser
	for _, ru := range repo.roomUsers {
		if ru.Room().Uid() == roomUid {
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
