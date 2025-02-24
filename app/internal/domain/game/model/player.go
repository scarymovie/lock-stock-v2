package model

import (
	"lock-stock-v2/internal/domain/room_user/model"
)

type Player struct {
	user    *model.RoomUser
	balance uint
	status  string
}

func NewPlayer(user *model.RoomUser, balance uint) *Player {
	return &Player{user: user, balance: balance}
}
