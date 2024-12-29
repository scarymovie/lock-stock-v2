//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	externalDomain "lock-stock-v2/external/domain"
	externalHandlers "lock-stock-v2/external/handlers"
	externalUsecase "lock-stock-v2/external/usecase"
	internalDomainRepository "lock-stock-v2/internal/domain/repository"
	internalHandlers "lock-stock-v2/internal/handlers"
	internalUsecase "lock-stock-v2/internal/usecase"
	"lock-stock-v2/router"
	"net/http"
)

// InitializeRouter связывает все зависимости и возвращает готовый http.Handler.
func InitializeRouter() (http.Handler, error) {
	wire.Build(
		// Handlers
		internalHandlers.NewJoinRoom,
		wire.Bind(new(externalHandlers.JoinRoom), new(*internalHandlers.JoinRoom)),

		// Usecase
		internalUsecase.NewJoinRoomUsecase,
		wire.Bind(new(externalUsecase.JoinRoom), new(*internalUsecase.JoinRoomUsecase)),

		// Domain - репозиторий покрывает два интерфейса
		internalDomainRepository.NewInMemoryRoomRepository,
		wire.Bind(new(externalDomain.RoomFinder), new(*internalDomainRepository.InMemoryRoomRepository)),
		//wire.Bind(new(externalDomain.RoomRepository), new(*internalDomainRepository.InMemoryRoomRepository)),

		// Роутер
		router.NewRouter,
	)

	return nil, nil
}
