package usecase

import "lock-stock-v2/external/domain"

type StartGameRequest struct {
	Room domain.Room
	User domain.User
}

type StartGame interface {
	StartGame(req StartGameRequest) error
}
