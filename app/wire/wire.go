//go:build wireinject
// +build wireinject

package wire

import (
	"errors"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
	externalTransactionManager "lock-stock-v2/external/transaction"
	externalWebSocket "lock-stock-v2/external/websocket"
	"lock-stock-v2/handlers/http/room"
	"lock-stock-v2/handlers/http/user"
	"lock-stock-v2/handlers/http/ws"
	gameRepository "lock-stock-v2/internal/domain/game/repository"
	gameService "lock-stock-v2/internal/domain/game/service"
	roomRepository "lock-stock-v2/internal/domain/room/repository"
	roomService "lock-stock-v2/internal/domain/room/service"
	roomUserRepository "lock-stock-v2/internal/domain/room_user/repository"
	roomUserService "lock-stock-v2/internal/domain/room_user/service"
	userRepository "lock-stock-v2/internal/domain/user/repository"
	userService "lock-stock-v2/internal/domain/user/service"
	internalPostgresRepository "lock-stock-v2/internal/infrastructure/postgres"
	internalPostgresTransactionManager "lock-stock-v2/internal/infrastructure/postgres/transaction"
	internalWebSocket "lock-stock-v2/internal/infrastructure/websocket"
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

func ProvideTransactionManager(pool *pgxpool.Pool) externalTransactionManager.TransactionManager {
	return internalPostgresTransactionManager.NewPostgresTransactionManager(pool)
}

func ProvideRoomHandler(
	joinRoomService *roomUserService.JoinRoomService,
	roomRepository roomRepository.RoomRepository,
	userRepository userRepository.UserRepository,
	roomUserService *roomUserService.RoomUserService,
	startGameService *roomService.StartGameService,
	createBetService *gameService.CreateBetService,
	sendAnswerService *gameService.SendAnswer,
	playerRepository gameRepository.PlayerRepository,
	roundRepository gameRepository.RoundRepository,
	betRepository gameRepository.BetRepository,
	gameRepository gameRepository.GameRepository,
	roundPlayerLogRepository gameRepository.RoundPlayerLogRepository,
	webSocket externalWebSocket.Manager,
	transactionManager externalTransactionManager.TransactionManager,
) *room.RoomHandler {
	return room.NewRoomHandler(
		joinRoomService,
		roomRepository,
		userRepository,
		roomUserService,
		startGameService,
		createBetService,
		playerRepository,
		roundRepository,
		betRepository,
		gameRepository,
		roundPlayerLogRepository,
		webSocket,
		sendAnswerService,
		transactionManager,
	)
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

func ProvideRoomUserRepository(db *pgxpool.Pool) roomUserRepository.RoomUserRepository {
	return internalPostgresRepository.NewPostgresRoomUserRepository(db)
}

func ProvideGameRepository(db *pgxpool.Pool) gameRepository.GameRepository {
	return internalPostgresRepository.NewPostgresGameRepository(db)
}

func ProvidePlayerRepository(db *pgxpool.Pool) gameRepository.PlayerRepository {
	return internalPostgresRepository.NewPostgresPlayerRepository(db)
}

func ProvideRoundRepository(db *pgxpool.Pool) gameRepository.RoundRepository {
	return internalPostgresRepository.NewPostgresRoundRepository(db)
}

func ProvideBetRepository(db *pgxpool.Pool) gameRepository.BetRepository {
	return internalPostgresRepository.NewPostgresBetRepository(db)
}

func ProvideRoundPlayerLogRepository(db *pgxpool.Pool) gameRepository.RoundPlayerLogRepository {
	return internalPostgresRepository.NewPostgresRoundPlayerLogRepository(db)
}

func InitializeRouter() (http.Handler, error) {
	wire.Build(
		// Подключение к PostgreSQL
		ProvidePostgresPool,
		ProvideTransactionManager,

		// Repositories
		ProvideRoomRepository,
		ProvideUserRepository,
		ProvideRoomUserRepository,
		ProvideGameRepository,
		ProvidePlayerRepository,
		ProvideRoundRepository,
		ProvideRoundPlayerLogRepository,
		ProvideBetRepository,

		// Services
		roomUserService.NewJoinRoomService,
		roomUserService.NewRoomUserService,
		userService.NewCreateUser,
		roomService.NewStartGameService,
		gameService.NewCreateGameService,
		gameService.NewCreateBetService,
		gameService.NewCreateRoundService,
		gameService.NewCreateRoundPlayerLog,
		gameService.NewSendAnswer,
		gameService.NewRoundObserver,

		// Handlers
		ProvideWebSocketHandler,
		ProvideUserHandler,
		ProvideRoomHandler,

		// WebSocket
		ProvideWebSocketManager,

		// Роутер
		router.NewRouter,
	)

	return nil, nil
}
