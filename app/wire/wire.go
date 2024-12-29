//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	externalHandlers "lock-stock-v2/external/handlers"
	internalHandlers "lock-stock-v2/internal/handlers"
	"lock-stock-v2/router"
	"net/http"
)

// InitializeRouter связывает все зависимости и возвращает готовый http.Handler.
func InitializeRouter() (http.Handler, error) {
	wire.Build(
		// Handlers
		internalHandlers.NewJoinRoom,
		wire.Bind(new(externalHandlers.JoinRoom), new(*internalHandlers.JoinRoom)), // Привязка интерфейса к реализации
		// Роутер
		router.NewRouter, // Роутер зависит от JoinRoom
	)

	return nil, nil
}
