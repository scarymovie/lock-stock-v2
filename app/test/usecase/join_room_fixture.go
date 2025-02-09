package usecase

import (
	model2 "lock-stock-v2/internal/domain/room/model"
	"lock-stock-v2/internal/domain/user/model"
)

func FakeUser(uid string, name string) *model.User {
	return model.NewUser(uid, name)
}

func FakeRoom(uid string) *model2.Room {
	return model2.NewRoom(uid, "StatusPending")
}
