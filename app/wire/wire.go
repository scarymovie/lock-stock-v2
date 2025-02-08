//go:build wireinject
// +build wireinject

package wire

import (
	"errors"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
	externalDomain "lock-stock-v2/external/domain"
	externalUsecase "lock-stock-v2/external/usecase"
	externalWebSocket "lock-stock-v2/external/websocket"
	"lock-stock-v2/handlers"
	internalDomainService "lock-stock-v2/internal/domain/service"
	internalPostgresRepository "lock-stock-v2/internal/infrastructure/postgres"
	internalWebSocket "lock-stock-v2/internal/infrastructure/websocket"
	internalUsecase "lock-stock-v2/internal/usecase"
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

func InitializeRouter() (http.Handler, error) {
	wire.Build(
		// Подключение к PostgreSQL
		ProvidePostgresPool,

		// Services
		internalDomainService.NewRoomService,

		// Handlers
		handlers.NewJoinRoom,
		handlers.NewGetRooms,
		handlers.NewCreateUser,
		handlers.NewStartGame,
		handlers.NewWebSocketHandler,

		// WebSocket
		ProvideWebSocketManager,

		// Usecase
		internalUsecase.NewJoinRoomUsecase,
		wire.Bind(new(externalUsecase.JoinRoom), new(*internalUsecase.JoinRoomUsecase)),

		internalUsecase.NewCreateUser,
		wire.Bind(new(externalUsecase.CreateUser), new(*internalUsecase.CreateUser)),

		internalUsecase.NewStartGameUsecase,
		wire.Bind(new(externalUsecase.StartGame), new(*internalUsecase.StartGameUsecase)),

		// Domain
		internalPostgresRepository.NewPostgresRoomRepository,
		wire.Bind(new(externalDomain.RoomFinder), new(*internalPostgresRepository.RoomRepository)),
		wire.Bind(new(externalDomain.RoomRepository), new(*internalPostgresRepository.RoomRepository)),

		internalPostgresRepository.NewPostgresUserRepository,
		wire.Bind(new(externalDomain.UserFinder), new(*internalPostgresRepository.UserRepository)),
		wire.Bind(new(externalDomain.UserRepository), new(*internalPostgresRepository.UserRepository)),

		internalPostgresRepository.NewPostgresRoomUserRepository,
		wire.Bind(new(externalDomain.RoomUserRepository), new(*internalPostgresRepository.RoomUserRepository)),

		// Роутер
		router.NewRouter,
	)

	return nil, nil
}
