package model

import (
	userModel "lock-stock-v2/internal/domain/user/model"
)

type PlayerStatus string

const (
	StatusLost    PlayerStatus = "lost"
	StatusPlaying PlayerStatus = "playing"
)

type Player struct {
	user    *userModel.User
	balance int
	status  PlayerStatus
	game    *LockStockGame
}

func NewPlayer(user *userModel.User, balance int, status PlayerStatus, game *LockStockGame) *Player {
	return &Player{user: user, balance: balance, status: status, game: game}
}

func (p *Player) User() *userModel.User {
	return p.user
}

func (p *Player) SetUser(user *userModel.User) {
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

func (p *Player) SetBalance(balance int) {
	p.balance = balance
}

func (p *Player) Balance() int {
	return p.balance
}
