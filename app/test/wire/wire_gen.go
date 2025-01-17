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

func InitializeTestCreateUserResult() TestCreateUserResult {
	userRepository := ProvideInMemoryUserRepository()
	createUser := usecase.NewCreateUser(userRepository)
	testCreateUserResult := TestCreateUserResult{
		CreateUser: createUser,
		UserRepo:   userRepository,
	}
	return testCreateUserResult
}

// wire.go:

type TestJoinRoomResult struct {
	JoinRoom usecase2.JoinRoom

	RoomRepo     *inMemory.RoomRepository
	UserRepo     *inMemory.UserRepository
	RoomUserRepo *inMemory.RoomUserRepository
	WsManager    websocket.Manager
}

type TestCreateUserResult struct {
	CreateUser usecase2.CreateUser
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

func ProvideInMemoryWebSocketManager() websocket.Manager {
	return inMemory.NewInMemoryWebSocketManager()
}

var testSetWithStruct = wire.NewSet(

	ProvideInMemoryRoomRepository, wire.Bind(new(domain.RoomRepository), new(*inMemory.RoomRepository)), ProvideInMemoryUserRepository, wire.Bind(new(domain.UserRepository), new(*inMemory.UserRepository)), ProvideInMemoryRoomUserRepository, wire.Bind(new(domain.RoomUserRepository), new(*inMemory.RoomUserRepository)), ProvideInMemoryWebSocketManager, usecase.NewJoinRoomUsecase, wire.Bind(new(usecase2.JoinRoom), new(*usecase.JoinRoomUsecase)), wire.Struct(new(TestJoinRoomResult), "*"),
)

var testSetCreateUser = wire.NewSet(

	ProvideInMemoryUserRepository, wire.Bind(new(domain.UserRepository), new(*inMemory.UserRepository)), usecase.NewCreateUser, wire.Bind(new(usecase2.CreateUser), new(*usecase.CreateUser)), wire.Struct(new(TestCreateUserResult), "*"),
)
