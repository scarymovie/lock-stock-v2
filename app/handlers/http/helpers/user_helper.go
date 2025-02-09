package helpers

import (
	"fmt"
	"lock-stock-v2/internal/domain/user/model"
	userRepository "lock-stock-v2/internal/domain/user/repository"
	"lock-stock-v2/middleware"
	"net/http"
)

func GetUserFromRequest(r *http.Request) (*model.User, error) {
	user, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		return nil, &UserNotFoundError{Code: http.StatusUnauthorized, Message: "Unauthorized"}
	}
	return user, nil
}

func GetUserFromString(userId string, repository userRepository.UserRepository) (*model.User, error) {
	user, err := repository.FindById(userId)
	if err != nil {
		_ = fmt.Errorf("user %s not found", userId)
		return nil, &UserNotFoundError{Code: http.StatusUnauthorized, Message: "Unauthorized"}
	}
	return user, nil
}

func (e *UserNotFoundError) Error() string {
	return e.Message
}

type UserNotFoundError struct {
	Code    int
	Message string
}
