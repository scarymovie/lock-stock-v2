package usecase

import (
	"lock-stock-v2/internal/domain"
)

func FakeUser(uid string) *domain.User {
	return &domain.User{
		Uid: uid,
		Id:  1,
	}
}

func FakeRoom(rid string) *domain.Room {
	return &domain.Room{
		Uid: rid,
		Id:  1,
	}
}
