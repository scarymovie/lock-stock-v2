package helpers

import (
	"lock-stock-v2/external/domain"
	"lock-stock-v2/middleware"
	"net/http"
)

func GetUserFromRequest(r *http.Request) (domain.User, error) {
	user, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
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
