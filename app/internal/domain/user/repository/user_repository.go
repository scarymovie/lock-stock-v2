package repository

import "lock-stock-v2/internal/domain/user/model"

type UserRepository interface {
	SaveUser(user *model.User) error
	FindById(userId string) (*model.User, error)
}
