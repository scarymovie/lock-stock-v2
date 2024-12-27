package wire

import (
	"github.com/google/wire"
	"lock-stock-v2/internal/handlers"
	"lock-stock-v2/router"
	"net/http"
)

// InitializeRouter связывает все зависимости и возвращает готовый http.Handler.
func InitializeRouter() (http.Handler, error) {
	wire.Build(
		// Handlers
		handlers.NewJoinRoomHandler,
		// Роутер
		router.NewRouter, // Роутер зависит от JoinRoom
	)

	return nil, nil
}
