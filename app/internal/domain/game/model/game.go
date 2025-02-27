package model

import (
	"lock-stock-v2/internal/domain/room/model"
	"time"
)

type LockStockGame struct {
	uid              string
	actionDuration   string
	questionDuration string
	room             *model.Room
	createdAt        time.Time
}

func NewLockStockGame(uid string, actionDuration string, questionDuration string, room *model.Room) *LockStockGame {
	return &LockStockGame{uid: uid, actionDuration: actionDuration, questionDuration: questionDuration, room: room}
}

func (l *LockStockGame) Uid() string {
	return l.uid
}

func (l *LockStockGame) SetUid(uid string) {
	l.uid = uid
}

func (l *LockStockGame) SetActionDuration(actionDuration string) {
	l.actionDuration = actionDuration
}

func (l *LockStockGame) SetQuestionDuration(questionDuration string) {
	l.questionDuration = questionDuration
}

func (l *LockStockGame) Room() *model.Room {
	return l.room
}

func (l *LockStockGame) SetRoom(room *model.Room) {
	l.room = room
}

func (l *LockStockGame) CreatedAt() time.Time {
	return l.createdAt
}

func (l *LockStockGame) SetCreatedAt(createdAt time.Time) {
	l.createdAt = createdAt
}

func (l *LockStockGame) ActionDuration() string {
	return l.actionDuration
}

func (l *LockStockGame) QuestionDuration() string {
	return l.questionDuration
}
