package middleware

import (
	"context"
	"errors"
	userModel "lock-stock-v2/internal/domain/user/model"
	"lock-stock-v2/internal/domain/user/repository"
	"net/http"
)

type UserContextKey string

const UserKey UserContextKey = "user"

func UserAuthMiddleware(userRepository repository.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userId := r.Header.Get("Authorization")
			if userId == "" {
				http.Error(w, "Authorization header missing", http.StatusUnauthorized)
				return
			}

			user, err := userRepository.FindById(userId)
			if err != nil {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserFromContext(ctx context.Context) (*userModel.User, error) {
	user, ok := ctx.Value(UserKey).(*userModel.User)
	if !ok {
		return nil, errors.New("user not found in context")
	}
	return user, nil
}
