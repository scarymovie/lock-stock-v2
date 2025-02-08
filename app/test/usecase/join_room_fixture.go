package usecase

import (
	"lock-stock-v2/internal/domain"
)

func FakeUser(uid string, name string) *domain.User {
	return domain.NewUser(uid, name)
}

func FakeRoom(uid string) *domain.Room {
	return domain.NewRoom(uid, "StatusPending")
}
