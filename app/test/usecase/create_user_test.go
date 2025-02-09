package usecase

import (
	"github.com/stretchr/testify/require"
	"lock-stock-v2/internal/domain/user/service"
	wireTest "lock-stock-v2/test/wire"
	"testing"
)

func TestCreateUser_Success(t *testing.T) {
	deps := wireTest.InitializeTestCreateUserResult()

	raw := service.RawCreateUser{Name: "Alice"}
	user, err := deps.CreateUser.Do(raw)

	require.NoError(t, err)
	require.NotNil(t, user)

	require.Equal(t, 1, deps.UserRepo.Count(), "ожидается, что 1 пользователь сохранён")
	require.Equal(t, "Alice", user.Name())
}
