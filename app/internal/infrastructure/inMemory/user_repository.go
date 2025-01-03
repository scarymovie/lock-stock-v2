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
		users: map[string]*domain.User{
			"user1": {Id: "user1"},
			"user2": {Id: "user2"},
			"user3": {Id: "user3"},
		},
	}
}

func (repo *UserRepository) FindById(userId string) (api.User, error) {
	user, exists := repo.users[userId]
	if !exists {
		return nil, errors.New("room not found")
	}
	return user, nil
}

func (repo *UserRepository) Save(user api.User) error {
	u, ok := user.(*domain.User)
	if !ok {
		return errors.New("invalid room type")
	}

	repo.users[u.GetUserId()] = u
	return nil
}
