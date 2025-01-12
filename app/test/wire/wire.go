//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"

	externalDomain "lock-stock-v2/external/domain"
	externalUsecase "lock-stock-v2/external/usecase"
	externalWs "lock-stock-v2/external/websocket"

	internalUsecase "lock-stock-v2/internal/usecase"

	"lock-stock-v2/test/inMemory"
)

type TestJoinRoomResult struct {
	JoinRoom externalUsecase.JoinRoom

	RoomRepo     *inMemory.RoomRepository
	UserRepo     *inMemory.UserRepository
	RoomUserRepo *inMemory.RoomUserRepository
	WsManager    externalWs.Manager
}

func ProvideInMemoryRoomRepository() *inMemory.RoomRepository {
	return inMemory.NewInMemoryRoomRepository()
}

func ProvideInMemoryUserRepository() *inMemory.UserRepository {
	return inMemory.NewInMemoryUserRepository()
}

func ProvideInMemoryRoomUserRepository() *inMemory.RoomUserRepository {
	return inMemory.NewInMemoryRoomUserRepository()
}

func ProvideInMemoryWebSocketManager() externalWs.Manager {
	return inMemory.NewInMemoryWebSocketManager()
}

var testSetWithStruct = wire.NewSet(

	ProvideInMemoryRoomRepository,
	wire.Bind(new(externalDomain.RoomRepository), new(*inMemory.RoomRepository)),

	ProvideInMemoryUserRepository,
	wire.Bind(new(externalDomain.UserRepository), new(*inMemory.UserRepository)),

	ProvideInMemoryRoomUserRepository,
	wire.Bind(new(externalDomain.RoomUserRepository), new(*inMemory.RoomUserRepository)),

	ProvideInMemoryWebSocketManager,

	internalUsecase.NewJoinRoomUsecase,
	wire.Bind(new(externalUsecase.JoinRoom), new(*internalUsecase.JoinRoomUsecase)),

	wire.Struct(new(TestJoinRoomResult), "*"),
)

func InitializeTestJoinRoomResult() TestJoinRoomResult {
	wire.Build(testSetWithStruct)
	return TestJoinRoomResult{}
}
