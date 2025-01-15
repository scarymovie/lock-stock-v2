package usecase

import (
	"github.com/google/uuid"
	"lock-stock-v2/external/domain"
	"lock-stock-v2/external/usecase"
	internalDomain "lock-stock-v2/internal/domain"
)

type CreateUser struct {
	userRepository domain.UserRepository
}

func NewCreateUser(userRepository domain.UserRepository) *CreateUser {
	return &CreateUser{userRepository: userRepository}
}

func (cu *CreateUser) Do(RawUser usecase.RawCreateUser) (domain.User, error) {
	newUser := internalDomain.User{
		Uid:  "user-" + uuid.New().String(),
		Name: RawUser.Name,
	}

	if err := cu.userRepository.SaveUser(&newUser); err != nil {
		return nil, err
	}

	return &newUser, nil
}
