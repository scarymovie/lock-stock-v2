package service

import (
	gameService "lock-stock-v2/internal/domain/game/service"
	roomModel "lock-stock-v2/internal/domain/room/model"
	"lock-stock-v2/internal/domain/room/repository"
	roomUserRepository "lock-stock-v2/internal/domain/room_user/repository"
	userModel "lock-stock-v2/internal/domain/user/model"
	"log"
)

type StartGameService struct {
	roomRepository     repository.RoomRepository
	roomUserRepository roomUserRepository.RoomUserRepository
	createGame         *gameService.CreateGameService
}

type StartGameRequest struct {
	Room *roomModel.Room
	User *userModel.User
}

func NewStartGameService(roomRepository repository.RoomRepository, roomUserRepository roomUserRepository.RoomUserRepository, createGame *gameService.CreateGameService) *StartGameService {
	return &StartGameService{roomRepository: roomRepository, roomUserRepository: roomUserRepository, createGame: createGame}
}

func (uc *StartGameService) StartGame(req StartGameRequest) error {

	req.Room.SetStatus(roomModel.StatusStarted)
	if err := uc.roomRepository.UpdateRoomStatus(req.Room); err != nil {
		log.Println("Failed to update room status:", err)
		return err
	}

	_, err := uc.createGame.CreateGame(req.Room)
	if err != nil {
		return err
	}
	return nil
}
