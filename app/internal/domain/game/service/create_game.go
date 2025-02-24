package service

import (
	"github.com/google/uuid"
	"lock-stock-v2/internal/domain/game/model"
	roomModel "lock-stock-v2/internal/domain/room/model"
	roomUserRepository "lock-stock-v2/internal/domain/room_user/repository"
	"log"
)

type CreateGameService struct {
	roomUserRepo roomUserRepository.RoomUserRepository
	roundService CreateRoundService
}

func NewCreateGameService(roomUserRepo roomUserRepository.RoomUserRepository) *CreateGameService {
	return &CreateGameService{roomUserRepo: roomUserRepo}
}

func (cgs *CreateGameService) CreateGame(room *roomModel.Room) *model.LockStockGame {

	roomUsers, err := cgs.roomUserRepo.FindByRoom(room)
	if err != nil {
		log.Println("aassdd")
	}

	var players []*model.Player
	for _, roomUser := range roomUsers {
		players = append(players, model.NewPlayer(roomUser, 5000))
	}

	game := model.NewLockStockGame("game"+uuid.New().String(), "30", "30", room, players)

	cgs.roundService.CreateRound(game, players)
	return game
}
