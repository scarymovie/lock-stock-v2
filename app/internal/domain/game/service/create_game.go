package service

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"lock-stock-v2/external/websocket"
	"lock-stock-v2/internal/domain/game/model"
	"lock-stock-v2/internal/domain/game/repository"
	roomModel "lock-stock-v2/internal/domain/room/model"
	roomUserRepository "lock-stock-v2/internal/domain/room_user/repository"
	userModel "lock-stock-v2/internal/domain/user/model"
	"log"
)

type CreateGameService struct {
	roomUserRepo     roomUserRepository.RoomUserRepository
	gameRepository   repository.GameRepository
	playerRepository repository.PlayerRepository
	roundService     *CreateRoundService
	webSocket        websocket.Manager
}

func NewCreateGameService(
	roomUserRepo roomUserRepository.RoomUserRepository,
	gameRepository repository.GameRepository,
	playerRepository repository.PlayerRepository,
	roundService *CreateRoundService,
	webSocket websocket.Manager) *CreateGameService {
	return &CreateGameService{
		roomUserRepo:     roomUserRepo,
		gameRepository:   gameRepository,
		playerRepository: playerRepository,
		roundService:     roundService,
		webSocket:        webSocket,
	}
}

func (cgs *CreateGameService) CreateGame(ctx context.Context, tx pgx.Tx, room *roomModel.Room) (*model.LockStockGame, error) {

	roomUsers, err := cgs.roomUserRepo.FindByRoom(room)
	if err != nil {
		log.Println("Error finding room users:", err)
		return nil, err
	}

	game := model.NewLockStockGame("game-"+uuid.New().String(), "30", "30", room)
	err = cgs.gameRepository.Save(ctx, tx, game)
	if err != nil {
		log.Println("Error saving game:", err)
		return nil, err
	}

	var playersFromRoomUsers []*userModel.User
	for _, roomUser := range roomUsers {
		playersFromRoomUsers = append(playersFromRoomUsers, roomUser.User())
	}

	players, playersData, err := cgs.createPlayers(ctx, tx, playersFromRoomUsers, game)
	if err != nil {
		log.Println("Error creating players:", err)
		return nil, err
	}

	err = cgs.publishGameStartedEvent(room, game, playersData)
	if err != nil {
		log.Println("Error publishing game started event:", err)
		return nil, err
	}

	err = cgs.roundService.CreateRound(ctx, tx, game, players)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (cgs *CreateGameService) createPlayers(ctx context.Context, tx pgx.Tx, users []*userModel.User, game *model.LockStockGame) ([]*model.Player, []map[string]interface{}, error) {
	var players []*model.Player
	var playersData []map[string]interface{}

	for _, user := range users {
		player := model.NewPlayer(user, 25000, model.StatusPlaying, game)
		err := cgs.playerRepository.Save(ctx, tx, player)
		if err != nil {
			log.Println("Error saving player:", err)
			return nil, nil, err
		}
		players = append(players, player)
		playersData = append(playersData, map[string]interface{}{
			"userId":  player.User().Uid(),
			"balance": player.Balance(),
		})
	}
	return players, playersData, nil
}

func (cgs *CreateGameService) publishGameStartedEvent(room *roomModel.Room, game *model.LockStockGame, playersData []map[string]interface{}) error {
	body := map[string]interface{}{
		"roomUid":          room.Uid(),
		"questionDuration": game.QuestionDuration(),
		"actionDuration":   game.ActionDuration(),
		"players":          playersData,
	}
	message := map[string]interface{}{
		"event": "start_game",
		"body":  body,
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal WebSocket message: %v\n", err)
		return err
	}

	log.Println(string(jsonMessage))
	cgs.webSocket.PublishToRoom(room.Uid(), jsonMessage)
	return nil
}
