package router

import (
	"github.com/go-chi/chi/v5"
	"lock-stock-v2/external/domain"
	"lock-stock-v2/external/handlers"
	"lock-stock-v2/middleware"
	"net/http"
)

// NewRouter принимает зависимости через Wire.
func NewRouter(
	joinRoom handlers.JoinRoom,
	wsHandler handlers.WebSocketHandler,
	userFinder domain.UserFinder,
	createUser handlers.CreateUser,
) http.Handler {
	r := chi.NewRouter()

	// Использование нескольких middleware.
	r.With(
		middleware.LoggingMiddleware,
		middleware.UserAuthMiddleware(userFinder),
	).Post("/join/room/{roomId}", joinRoom.ServeHTTP)

	r.With(
		middleware.LoggingMiddleware,
	).Post("/user/create", createUser.ServeHTTP)

	r.Handle("/ws", wsHandler)

	return r
}
