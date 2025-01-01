package middleware

import (
	"context"
	"errors"
	"net/http"

	"lock-stock-v2/external/domain"
)

type UserContextKey string

const UserKey UserContextKey = "user"

func UserAuthMiddleware(userFinder domain.UserFinder) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userId := r.Header.Get("Authorization")
			if userId == "" {
				http.Error(w, "Authorization header missing", http.StatusUnauthorized)
				return
			}

			user, err := userFinder.FindById(userId)
			if err != nil {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			// Сохранить пользователя в контексте
			ctx := context.WithValue(r.Context(), UserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserFromContext(ctx context.Context) (domain.User, error) {
	user, ok := ctx.Value(UserKey).(domain.User)
	if !ok {
		return nil, errors.New("user not found in context")
	}
	return user, nil
}
