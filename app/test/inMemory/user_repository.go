package inMemory

import (
	"errors"
	api "lock-stock-v2/external/domain"
	"lock-stock-v2/internal/domain"
)

type UserRepository struct {
	users map[string]*domain.User
}

func NewInMemoryUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]*domain.User),
	}
}

func (repo *UserRepository) FindById(userId string) (api.User, error) {
	user, exists := repo.users[userId]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (repo *UserRepository) SaveUser(user api.User) error {
	u, ok := user.(*domain.User)
	if !ok {
		return errors.New("invalid user type")
	}

	repo.users[u.GetUserUid()] = u
	return nil
}

func (repo *UserRepository) Count() int {
	return len(repo.users)
}
