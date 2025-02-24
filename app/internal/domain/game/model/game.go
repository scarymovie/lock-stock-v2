package model

import "lock-stock-v2/internal/domain/room/model"

type LockStockGame struct {
	uid              string
	actionDuration   string
	questionDuration string
	room             *model.Room
	players          []*Player
}

func NewLockStockGame(uid string, actionDuration string, questionDuration string, room *model.Room, players []*Player) *LockStockGame {
	return &LockStockGame{uid: uid, actionDuration: actionDuration, questionDuration: questionDuration, room: room, players: players}
}

func (l *LockStockGame) ActionDuration() string {
	return l.actionDuration
}

func (l *LockStockGame) QuestionDuration() string {
	return l.questionDuration
}
