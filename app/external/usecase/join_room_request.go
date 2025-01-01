package usecase

import "lock-stock-v2/external/domain"

type JoinRoomRequest struct {
	User domain.User
	Room domain.Room
}
