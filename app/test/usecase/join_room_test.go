package usecase

import (
	"github.com/stretchr/testify/require"
	services "lock-stock-v2/internal/domain/room_user/service"
	wireTest "lock-stock-v2/test/wire"
	"testing"
)

func TestJoinRoom_Success(t *testing.T) {

	deps := wireTest.InitializeTestJoinRoomResult()

	require.Equal(t, 0, deps.UserRepo.Count())
	require.Equal(t, 0, deps.RoomRepo.Count())
	require.Equal(t, 0, deps.RoomUserRepo.Count())

	user := FakeUser("user-123", "user-name")
	room := FakeRoom("room-456")

	deps.UserRepo.SaveUser(user)
	deps.RoomRepo.Save(room)

	request := services.JoinRoomRequest{
		User: user,
		Room: room,
	}

	err := deps.JoinRoom.JoinRoom(request)
	require.NoError(t, err, "ожидается, что при корректных данных не будет ошибки")

	require.Equal(t, 1, deps.UserRepo.Count(), "должен быть 1 пользователь")

	require.Equal(t, 1, deps.RoomRepo.Count(), "должна быть 1 комната")

	require.Equal(t, 1, deps.RoomUserRepo.Count(), "должна быть 1 запись room-user")

	// При желании проверяем, что вебсокет тоже получил 1 сообщение и т.д.
	// Если InMemoryWebSocketManager имеет метод GetMessagesForRoom, можно вызвать его:
	// msgs := deps.WsManager.(*inMemory.InMemoryWebSocketManager).GetMessagesForRoom("room-456")
	// require.Len(t, msgs, 1)
}
