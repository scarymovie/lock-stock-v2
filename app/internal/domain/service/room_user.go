package services

import (
	"lock-stock-v2/external/domain"
	"log"
)

type RoomUserService struct {
	roomUsersRepository domain.RoomUserRepository
}

func NewRoomService(roomUsersRepository domain.RoomUserRepository) *RoomUserService {
	return &RoomUserService{roomUsersRepository: roomUsersRepository}
}

func (s *RoomUserService) GetUsersByRoom(room domain.Room) ([]domain.RoomUser, error) {
	roomUsers, err := s.roomUsersRepository.FindByRoom(room)
	if err != nil {
		log.Println("Error fetching room users:", err)
		return nil, err
	}
	return roomUsers, nil
}

func (s *RoomUserService) IsUserInRoom(roomUsers []domain.RoomUser, user domain.User) bool {
	for _, ru := range roomUsers {
		if ru.GetUser().GetUserId() == user.GetUserId() {
			return true
		}
	}
	return false
}
