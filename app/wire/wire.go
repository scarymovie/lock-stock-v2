//go:build wireinject
// +build wireinject

package wire

import (
	"errors"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
	externalDomain "lock-stock-v2/external/domain"
	externalHandlers "lock-stock-v2/external/handlers"
	externalUsecase "lock-stock-v2/external/usecase"
	internalHandlers "lock-stock-v2/internal/handlers"
	internalPostgresRepository "lock-stock-v2/internal/infrastructure/postgres"
	internalUsecase "lock-stock-v2/internal/usecase"
	internalWebsocket "lock-stock-v2/internal/websocket"
	"lock-stock-v2/router"
	"net/http"
)

func ProvidePostgresPool() (*pgxpool.Pool, error) {
	// Получаем конфигурацию
	config := internalPostgresRepository.GetPostgresConfig()

	// Создаем подключение
	pool := internalPostgresRepository.NewPostgresConnection(config)
	if pool == nil {
		return nil, errors.New("failed to create postgres pool")
	}

	return pool, nil
}

func ProvideWebSocketManager() *internalWebsocket.WebSocketManager {
	manager := internalWebsocket.NewWebSocketManager()
	go manager.Run()
	return manager
}

// InitializeRouter связывает все зависимости и возвращает готовый http.Handler.
func InitializeRouter() (http.Handler, error) {
	wire.Build(
		// Подключение к PostgreSQL
		ProvidePostgresPool, // Провайдер для подключения

		// Handlers
		internalHandlers.NewJoinRoom,
		wire.Bind(new(externalHandlers.JoinRoom), new(*internalHandlers.JoinRoom)),

		internalHandlers.NewCreateUser,
		wire.Bind(new(externalHandlers.CreateUser), new(*internalHandlers.CreateUser)),

		ProvideWebSocketManager,
		internalHandlers.NewWebSocketHandler,
		wire.Bind(new(externalHandlers.WebSocketHandler), new(*internalHandlers.WebSocketHandler)),

		// Usecase
		internalUsecase.NewJoinRoomUsecase,
		wire.Bind(new(externalUsecase.JoinRoom), new(*internalUsecase.JoinRoomUsecase)),

		// Domain
		internalPostgresRepository.NewPostgresRoomRepository,
		wire.Bind(new(externalDomain.RoomFinder), new(*internalPostgresRepository.RoomRepository)),

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
