package model

import (
	"lock-stock-v2/internal/domain/room_user/model"
)

type PlayerStatus string

const (
	StatusLost    PlayerStatus = "lost"
	StatusPlaying PlayerStatus = "playing"
)

type Player struct {
	user    *model.RoomUser
	balance uint
	status  PlayerStatus
	game    *LockStockGame
}

func NewPlayer(user *model.RoomUser, balance uint, status PlayerStatus, game *LockStockGame) *Player {
	return &Player{user: user, balance: balance, status: status, game: game}
}
