package repository

import (
	"errors"
	api "lock-stock-v2/external/domain"
	"lock-stock-v2/internal/domain"
)

type InMemoryUserRepository struct {
	users map[string]*domain.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: map[string]*domain.User{
			"user1": {Id: "user1"},
			"user2": {Id: "user2"},
			"user3": {Id: "user3"},
		},
	}
}

func (repo *InMemoryUserRepository) FindById(userId string) (api.User, error) {
	user, exists := repo.users[userId]
	if !exists {
		return nil, errors.New("room not found")
	}
	return user, nil
}

func (repo *InMemoryUserRepository) Save(user api.User) error {
	u, ok := user.(*domain.User)
	if !ok {
		return errors.New("invalid room type")
	}

	repo.users[u.GetId()] = u
	return nil
}
