package router

import (
	"github.com/go-chi/chi/v5"
	"lock-stock-v2/external/handlers"
	"lock-stock-v2/middleware"
	"net/http"
)

// NewRouter принимает зависимости через Wire.
func NewRouter(joinRoom handlers.JoinRoom) http.Handler {
	r := chi.NewRouter()

	// Регистрация маршрута с middleware и обработчиком.
	r.With(middleware.LoggingMiddleware).Post("/room/{id}", joinRoom.ServeHTTP)

	return r
}
