package usecase

import "lock-stock-v2/external/domain"

type CreateUser interface {
	Do(user RawCreateUser) (domain.User, error)
}
