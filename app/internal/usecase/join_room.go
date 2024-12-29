package usecase

import (
	"fmt"
	"lock-stock-v2/external/domain"
	"lock-stock-v2/external/usecase"
	"log"
)

type JoinRoomUsecase struct {
	roomFinder domain.RoomFinder
}

func NewJoinRoom() *JoinRoomUsecase {
	return &JoinRoomUsecase{}
}

func (s JoinRoomUsecase) JoinRoom(request usecase.JoinRoomRequest) error {
	room, err := s.roomFinder.FindById(request.RoomId)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if room == nil {
		fmt.Printf("Room %v does not exist", request.RoomId)
	}
	fmt.Printf("usecase Player %s joined room %s\n", request.PlayerId, request.RoomId)
	return nil
}