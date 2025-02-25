package model

import (
	"github.com/google/uuid"
	userModel "lock-stock-v2/internal/domain/user/model"
)

type PlayerStatus string

const (
	StatusLost    PlayerStatus = "lost"
	StatusPlaying PlayerStatus = "playing"
)

type Player struct {
	uid     string
	user    *userModel.User
	balance uint
	status  PlayerStatus
	game    *LockStockGame
}

func (p *Player) User() *userModel.User {
	return p.user
}

func (p *Player) SetUser(user *userModel.User) {
	p.user = user
}

func NewPlayer(user *userModel.User, balance uint, status PlayerStatus, game *LockStockGame) *Player {
	return &Player{uid: "player-" + uuid.New().String(), user: user, balance: balance, status: status, game: game}
}

func (p *Player) Uid() string {
	return p.uid
}

func (p *Player) SetUid(uid string) {
	p.uid = uid
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
