//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"lock-stock-v2/internal/domain/room/repository"
	externalDomain "lock-stock-v2/internal/domain/room_user/repository"
	internalUsecase "lock-stock-v2/internal/domain/room_user/service"
	repository2 "lock-stock-v2/internal/domain/user/repository"
	"lock-stock-v2/internal/domain/user/service"
	externalWs "lock-stock-v2/internal/websocket"

	externalUsecase "lock-stock-v2/external/usecase"
	"lock-stock-v2/test/inMemory"
)

type TestJoinRoomResult struct {
	JoinRoom externalUsecase.JoinRoom

	RoomRepo     *inMemory.RoomRepository
	UserRepo     *inMemory.UserRepository
	RoomUserRepo *inMemory.RoomUserRepository
	WsManager    externalWs.Manager
}

type TestCreateUserResult struct {
	CreateUser externalUsecase.CreateUser
	UserRepo   *inMemory.UserRepository
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
	wire.Bind(new(repository.RoomRepository), new(*inMemory.RoomRepository)),

	ProvideInMemoryUserRepository,
	wire.Bind(new(repository2.UserRepository), new(*inMemory.UserRepository)),

	ProvideInMemoryRoomUserRepository,
	wire.Bind(new(externalDomain.RoomUserRepository), new(*inMemory.RoomUserRepository)),

	ProvideInMemoryWebSocketManager,

	internalUsecase.NewJoinRoom,
	wire.Bind(new(externalUsecase.JoinRoom), new(*internalUsecase.JoinRoomService)),

	wire.Struct(new(TestJoinRoomResult), "*"),
)

var testSetCreateUser = wire.NewSet(

	ProvideInMemoryUserRepository,
	wire.Bind(new(repository2.UserRepository), new(*inMemory.UserRepository)),

	service.NewCreateUser,
	wire.Bind(new(externalUsecase.CreateUser), new(*service.CreateUserService)),

	wire.Struct(new(TestCreateUserResult), "*"),
)

func InitializeTestJoinRoomResult() TestJoinRoomResult {
	wire.Build(testSetWithStruct)
	return TestJoinRoomResult{}
}

func InitializeTestCreateUserResult() TestCreateUserResult {
	wire.Build(testSetCreateUser)
	return TestCreateUserResult{}
}
