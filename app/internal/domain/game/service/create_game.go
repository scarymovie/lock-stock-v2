package service

import (
	"github.com/google/uuid"
	"lock-stock-v2/internal/domain/game/model"
	"lock-stock-v2/internal/domain/game/repository"
	roomModel "lock-stock-v2/internal/domain/room/model"
	roomUserRepository "lock-stock-v2/internal/domain/room_user/repository"
	"log"
)

type CreateGameService struct {
	roomUserRepo     roomUserRepository.RoomUserRepository
	gameRepository   repository.GameRepository
	playerRepository repository.PlayerRepository
	roundService     CreateRoundService
}

func NewCreateGameService(roomUserRepo roomUserRepository.RoomUserRepository) *CreateGameService {
	return &CreateGameService{roomUserRepo: roomUserRepo}
}

func (cgs *CreateGameService) CreateGame(room *roomModel.Room) *model.LockStockGame {

	roomUsers, err := cgs.roomUserRepo.FindByRoom(room)
	if err != nil {
		log.Println("aassdd")
	}

	game := model.NewLockStockGame("game"+uuid.New().String(), "30", "30", room)
	cgs.gameRepository.Save(game)

	var players []*model.Player
	for _, roomUser := range roomUsers {
		player := model.NewPlayer(roomUser, 25000, model.StatusPlaying, game)
		cgs.playerRepository.Save(player)
		players = append(players, player)
	}

	cgs.roundService.CreateRound(game, players)
	return game
}
