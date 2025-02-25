//go:build wireinject
// +build wireinject

package wire

import (
	"errors"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
	"lock-stock-v2/handlers/http/room"
	"lock-stock-v2/handlers/http/user"
	"lock-stock-v2/handlers/http/ws"
	gameRepository "lock-stock-v2/internal/domain/game/repository"
	roomRepository "lock-stock-v2/internal/domain/room/repository"
	roomService "lock-stock-v2/internal/domain/room/service"
	roomUserRepository "lock-stock-v2/internal/domain/room_user/repository"
	roomUserService "lock-stock-v2/internal/domain/room_user/service"
	userRepository "lock-stock-v2/internal/domain/user/repository"
	userService "lock-stock-v2/internal/domain/user/service"
	internalPostgresRepository "lock-stock-v2/internal/infrastructure/postgres"
	internalWebSocket "lock-stock-v2/internal/infrastructure/websocket"
	externalWebSocket "lock-stock-v2/internal/websocket"
	"lock-stock-v2/router"
	"net/http"
)

func ProvidePostgresPool() (*pgxpool.Pool, error) {
	config := internalPostgresRepository.GetPostgresConfig()
	pool := internalPostgresRepository.NewPostgresConnection(config)
	if pool == nil {
		return nil, errors.New("failed to create postgres pool")
	}

	return pool, nil
}

func ProvideWebSocketManager() externalWebSocket.Manager {
	manager := internalWebSocket.NewWebSocketManager()
	go manager.Run()
	return manager
}

func ProvideRoomHandler(
	joinRoomService *roomUserService.JoinRoomService,
	roomRepository roomRepository.RoomRepository,
	userRepository userRepository.UserRepository,
	roomUserService *roomUserService.RoomUserService,
	startGameService *roomService.StartGameService,
) *room.RoomHandler {
	return room.NewRoomHandler(joinRoomService, roomRepository, userRepository, roomUserService, startGameService)
}

func ProvideUserHandler(createUserService *userService.CreateUserService) *user.UserHandler {
	return user.NewUserHandler(createUserService)
}

func ProvideWebSocketHandler(manager externalWebSocket.Manager) *ws.WebSocketHandler {
	return ws.NewWebSocketHandler(manager)
}

func ProvideRoomRepository(db *pgxpool.Pool) roomRepository.RoomRepository {
	return internalPostgresRepository.NewPostgresRoomRepository(db)
}

func ProvideUserRepository(db *pgxpool.Pool) userRepository.UserRepository {
	return internalPostgresRepository.NewPostgresUserRepository(db)
}

func ProvideGameRepository(db *pgxpool.Pool) gameRepository.GameRepository {
	return internalPostgresRepository.NewPostgresGameRepository(db)
}

func ProvideRoomUserRepository(db *pgxpool.Pool) roomUserRepository.RoomUserRepository {
	return internalPostgresRepository.NewPostgresRoomUserRepository(db)
}

func InitializeRouter() (http.Handler, error) {
	wire.Build(
		// Подключение к PostgreSQL
		ProvidePostgresPool,

		// Services
		roomUserService.NewJoinRoomService,
		roomUserService.NewRoomUserService,
		userService.NewCreateUser,
		roomService.NewStartGameService,
		//gameService.NewCreateGameService,
		//gameService.NewCreateRoundService,

		// Handlers
		ProvideWebSocketHandler,
		ProvideUserHandler,
		ProvideRoomHandler,

		// WebSocket
		ProvideWebSocketManager,

		// Repositories
		ProvideRoomRepository,
		ProvideUserRepository,
		ProvideRoomUserRepository,
		//ProvideGameRepository,

		// Роутер
		router.NewRouter,
	)

	return nil, nil
}
