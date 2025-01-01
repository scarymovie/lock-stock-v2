package usecase

import (
	"lock-stock-v2/external/domain"
	"lock-stock-v2/external/usecase"
	internalDomain "lock-stock-v2/internal/domain"
	"log"
)

type JoinRoomUsecase struct {
	roomUserRepository domain.RoomUserRepository
}

func NewJoinRoomUsecase(roomUserRepository domain.RoomUserRepository) *JoinRoomUsecase {
	return &JoinRoomUsecase{
		roomUserRepository: roomUserRepository,
	}
}

func (s JoinRoomUsecase) JoinRoom(request usecase.JoinRoomRequest) error {
	roomUser := internalDomain.RoomUser{}
	roomUser.SetRoom(request.Room)
	roomUser.SetUser(request.User)
	err := s.roomUserRepository.Save(&roomUser)
	if err != nil {
		return err
	}

	log.Printf("usecase Player %s joined room %s\n", request.User.GetUserId(), request.Room.GetRoomId())
	return nil
}
