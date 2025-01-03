// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"lock-stock-v2/internal/handlers"
	"lock-stock-v2/internal/infrastructure/postgres"
	"lock-stock-v2/internal/usecase"
	"lock-stock-v2/internal/websocket"
	"lock-stock-v2/router"
	"net/http"
)

// Injectors from wire.go:

// InitializeRouter связывает все зависимости и возвращает готовый http.Handler.
func InitializeRouter() (http.Handler, error) {
	pool, err := ProvidePostgresPool()
	if err != nil {
		return nil, err
	}
	roomUserRepository := postgres.NewPostgresRoomUserRepository(pool)
	webSocketManager := ProvideWebSocketManager()
	joinRoomUsecase := usecase.NewJoinRoomUsecase(roomUserRepository, webSocketManager)
	roomRepository := postgres.NewPostgresRoomRepository(pool)
	joinRoom := handlers.NewJoinRoom(joinRoomUsecase, roomRepository)
	webSocketHandler := handlers.NewWebSocketHandler(webSocketManager)
	userRepository := postgres.NewPostgresUserRepository(pool)
	createUser := handlers.NewCreateUser(userRepository)
	handler := router.NewRouter(joinRoom, webSocketHandler, userRepository, createUser)
	return handler, nil
}

// wire.go:

func ProvidePostgresPool() (*pgxpool.Pool, error) {

	config := postgres.GetPostgresConfig()

	pool := postgres.NewPostgresConnection(config)
	if pool == nil {
		return nil, errors.New("failed to create postgres pool")
	}

	return pool, nil
}

func ProvideWebSocketManager() *websocket.WebSocketManager {
	manager := websocket.NewWebSocketManager()
	go manager.Run()
	return manager
}
