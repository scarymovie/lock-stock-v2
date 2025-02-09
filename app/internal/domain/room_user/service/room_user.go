package services

import (
	roomModel "lock-stock-v2/internal/domain/room/model"
	roomUserModel "lock-stock-v2/internal/domain/room_user/model"
	"lock-stock-v2/internal/domain/room_user/repository"
	userModel "lock-stock-v2/internal/domain/user/model"
	"log"
)

type RoomUserService struct {
	roomUsersRepository repository.RoomUserRepository
}

func NewRoomUserService(roomUsersRepository repository.RoomUserRepository) *RoomUserService {
	return &RoomUserService{roomUsersRepository: roomUsersRepository}
}

func (s *RoomUserService) GetUsersByRoom(room *roomModel.Room) ([]*roomUserModel.RoomUser, error) {
	roomUsers, err := s.roomUsersRepository.FindByRoom(room)
	if err != nil {
		log.Println("Error fetching room users:", err)
		return nil, err
	}
	return roomUsers, nil
}

func (s *RoomUserService) IsUserInRoom(roomUsers []*roomUserModel.RoomUser, user *userModel.User) bool {
	for _, ru := range roomUsers {
		u := ru.User()
		if u.Uid() == user.Uid() {
			return true
		}
	}
	return false
}
