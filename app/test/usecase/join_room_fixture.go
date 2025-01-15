package usecase

import (
	"lock-stock-v2/internal/domain"
)

func FakeUser(uid string, name string) *domain.User {
	return &domain.User{
		Uid:  uid,
		Id:   1,
		Name: name,
	}
}

func FakeRoom(rid string) *domain.Room {
	return &domain.Room{
		Uid: rid,
		Id:  1,
	}
}
