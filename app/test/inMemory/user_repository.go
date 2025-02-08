package inMemory

import (
	"errors"
	api "lock-stock-v2/external/domain"
	"lock-stock-v2/internal/domain/user/model"
)

type UserRepository struct {
	users map[string]*model.User
}

func NewInMemoryUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]*model.User),
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
	u, ok := user.(*model.User)
	if !ok {
		return errors.New("invalid user type")
	}

	repo.users[u.Uid()] = u
	return nil
}

func (repo *UserRepository) Count() int {
	return len(repo.users)
}
