package inMemory

import (
	"errors"
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

func (repo *UserRepository) FindById(userId string) (*model.User, error) {
	user, exists := repo.users[userId]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (repo *UserRepository) SaveUser(user *model.User) error {
	repo.users[user.Uid()] = user
	return nil
}

func (repo *UserRepository) Count() int {
	return len(repo.users)
}
