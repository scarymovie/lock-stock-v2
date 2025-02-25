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

func (p *Player) RoomUser() *model.RoomUser {
	return p.user
}

func (p *Player) SetUser(user *model.RoomUser) {
	p.user = user
}

func (p *Player) Status() PlayerStatus {
	return p.status
}

func (p *Player) SetStatus(status PlayerStatus) {
	p.status = status
}

func (p *Player) Game() *LockStockGame {
	return p.game
}

func (p *Player) SetGame(game *LockStockGame) {
	p.game = game
}

func (p *Player) SetBalance(balance uint) {
	p.balance = balance
}

func (p *Player) Balance() uint {
	return p.balance
}
