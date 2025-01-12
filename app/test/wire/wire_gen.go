// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"github.com/google/wire"
	"lock-stock-v2/external/domain"
	usecase2 "lock-stock-v2/external/usecase"
	"lock-stock-v2/external/websocket"
	"lock-stock-v2/internal/usecase"
	"lock-stock-v2/test/inMemory"
)

// Injectors from wire.go:

// Функция, которую Wire «замкнёт»:
func InitializeTestJoinRoomResult() TestJoinRoomResult {
	roomUserRepository := ProvideInMemoryRoomUserRepository()
	manager := ProvideInMemoryWebSocketManager()
	joinRoomUsecase := usecase.NewJoinRoomUsecase(roomUserRepository, manager)
	roomRepository := ProvideInMemoryRoomRepository()
	userRepository := ProvideInMemoryUserRepository()
	testJoinRoomResult := TestJoinRoomResult{
		JoinRoom:     joinRoomUsecase,
		RoomRepo:     roomRepository,
		UserRepo:     userRepository,
		RoomUserRepo: roomUserRepository,
		WsManager:    manager,
	}
	return testJoinRoomResult
}

// wire.go:

// Структура, куда мы хотим «сложить» все зависимости.
type TestJoinRoomResult struct {
	JoinRoom usecase2.JoinRoom
	// Ниже — именно *inMemory.RoomRepository, чтобы вы могли делать .Count() и т.д.
	RoomRepo     *inMemory.RoomRepository
	UserRepo     *inMemory.UserRepository
	RoomUserRepo *inMemory.RoomUserRepository
	WsManager    websocket.Manager
}

// --- Провайдеры, возвращающие конкретные структуры ---
func ProvideInMemoryRoomRepository() *inMemory.RoomRepository {
	return inMemory.NewInMemoryRoomRepository()
}

func ProvideInMemoryUserRepository() *inMemory.UserRepository {
	return inMemory.NewInMemoryUserRepository()
}

func ProvideInMemoryRoomUserRepository() *inMemory.RoomUserRepository {
	return inMemory.NewInMemoryRoomUserRepository()
}

func ProvideInMemoryWebSocketManager() websocket.Manager {
	return inMemory.NewInMemoryWebSocketManager()
}

// Набор wire
var testSetWithStruct = wire.NewSet(

	ProvideInMemoryRoomRepository, wire.Bind(new(domain.RoomRepository), new(*inMemory.RoomRepository)), ProvideInMemoryUserRepository, wire.Bind(new(domain.UserRepository), new(*inMemory.UserRepository)), ProvideInMemoryRoomUserRepository, wire.Bind(new(domain.RoomUserRepository), new(*inMemory.RoomUserRepository)), ProvideInMemoryWebSocketManager, usecase.NewJoinRoomUsecase, wire.Bind(new(usecase2.JoinRoom), new(*usecase.JoinRoomUsecase)), wire.Struct(new(TestJoinRoomResult), "*"),
)