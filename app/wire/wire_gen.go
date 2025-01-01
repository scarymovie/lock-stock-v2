// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"lock-stock-v2/internal/handlers"
	"lock-stock-v2/internal/infrastructure/inMemory"
	"lock-stock-v2/internal/usecase"
	"lock-stock-v2/router"
	"net/http"
)

// Injectors from wire.go:

// InitializeRouter связывает все зависимости и возвращает готовый http.Handler.
func InitializeRouter() (http.Handler, error) {
	roomUserRepository := inMemory.NewInMemoryRoomUserRepository()
	joinRoomUsecase := usecase.NewJoinRoomUsecase(roomUserRepository)
	roomRepository := inMemory.NewInMemoryRoomRepository()
	joinRoom := handlers.NewJoinRoom(joinRoomUsecase, roomRepository)
	userRepository := inMemory.NewInMemoryUserRepository()
	handler := router.NewRouter(joinRoom, userRepository)
	return handler, nil
}
