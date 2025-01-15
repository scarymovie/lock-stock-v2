package usecase

import (
	"github.com/stretchr/testify/require"
	externalUsecase "lock-stock-v2/external/usecase"
	internalDomain "lock-stock-v2/internal/domain"
	wireTest "lock-stock-v2/test/wire"
	"testing"
)

func TestCreateUser_Success(t *testing.T) {
	deps := wireTest.InitializeTestCreateUserResult()

	raw := externalUsecase.RawCreateUser{Name: "Alice"}
	user, err := deps.CreateUser.Do(raw)

	require.NoError(t, err)
	require.NotNil(t, user)

	require.Equal(t, 1, deps.UserRepo.Count(), "ожидается, что 1 пользователь сохранён")

	u, ok := user.(*internalDomain.User)
	require.True(t, ok, "должен быть *internalDomain.User")
	require.Equal(t, "Alice", u.Name)
}
