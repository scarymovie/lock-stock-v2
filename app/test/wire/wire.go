//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	externalWs "lock-stock-v2/external/websocket"
	"lock-stock-v2/internal/domain/room/repository"
	roomUserRepository "lock-stock-v2/internal/domain/room_user/repository"
	roomUserService "lock-stock-v2/internal/domain/room_user/service"
	userRepository "lock-stock-v2/internal/domain/user/repository"
	userService "lock-stock-v2/internal/domain/user/service"
	"lock-stock-v2/test/inMemory"
)

type TestJoinRoomResult struct {
	JoinRoom *roomUserService.JoinRoomService

	RoomRepo     *inMemory.RoomRepository
	UserRepo     *inMemory.UserRepository
	RoomUserRepo *inMemory.RoomUserRepository
	WsManager    externalWs.Manager
}

type TestCreateUserResult struct {
	CreateUser *userService.CreateUserService
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

func ProvideJoinRoomService(roomUserRepository roomUserRepository.RoomUserRepository,
	wsManager externalWs.Manager,
) *roomUserService.JoinRoomService {
	return roomUserService.NewJoinRoomService(roomUserRepository, wsManager)
}

func ProvideCreateUserService(userRepository userRepository.UserRepository) *userService.CreateUserService {
	return userService.NewCreateUser(userRepository)
}

func ProvideInMemoryWebSocketManager() externalWs.Manager {
	return inMemory.NewInMemoryWebSocketManager()
}

var testSetWithStruct = wire.NewSet(

	ProvideInMemoryRoomRepository,
	wire.Bind(new(repository.RoomRepository), new(*inMemory.RoomRepository)),

	ProvideInMemoryUserRepository,
	wire.Bind(new(userRepository.UserRepository), new(*inMemory.UserRepository)),

	ProvideInMemoryRoomUserRepository,
	wire.Bind(new(roomUserRepository.RoomUserRepository), new(*inMemory.RoomUserRepository)),

	ProvideInMemoryWebSocketManager,
	ProvideJoinRoomService,

	wire.Struct(new(TestJoinRoomResult), "*"),
)

var testSetCreateUser = wire.NewSet(

	ProvideInMemoryUserRepository,
	wire.Bind(new(userRepository.UserRepository), new(*inMemory.UserRepository)),

	ProvideCreateUserService,

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
