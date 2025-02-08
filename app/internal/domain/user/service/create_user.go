package service

import (
	"github.com/google/uuid"
	userModel "lock-stock-v2/internal/domain/user/model"
	"lock-stock-v2/internal/domain/user/repository"
)

type CreateUserService struct {
	userRepository repository.UserRepository
}

type RawCreateUser struct {
	Name string
}

func NewCreateUser(userRepository repository.UserRepository) *CreateUserService {
	return &CreateUserService{userRepository: userRepository}
}

func (cu *CreateUserService) Do(RawUser RawCreateUser) (*userModel.User, error) {
	newUser := userModel.NewUser(
		"user-"+uuid.New().String(),
		RawUser.Name,
	)

	if err := cu.userRepository.SaveUser(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}
