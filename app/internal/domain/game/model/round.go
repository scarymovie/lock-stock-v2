package model

import "github.com/google/uuid"

type Round struct {
	uid      string
	number   *uint
	question *Question
	buyIn    uint
	pot      uint
	game     *LockStockGame
}

func NewRound(number *uint, question *Question, buyIn uint, pot uint, game *LockStockGame) *Round {
	return &Round{uid: "round-" + uuid.New().String(), number: number, question: question, buyIn: buyIn, pot: pot, game: game}
}

func (r *Round) Uid() string {
	return r.uid
}

func (r *Round) SetUid(uid string) {
	r.uid = uid
}

func (r *Round) Number() *uint {
	return r.number
}

func (r *Round) SetNumber(number *uint) {
	r.number = number
}

func (r *Round) Question() *Question {
	return r.question
}

func (r *Round) SetQuestion(question *Question) {
	r.question = question
}

func (r *Round) BuyIn() uint {
	return r.buyIn
}

func (r *Round) SetBuyIn(buyIn uint) {
	r.buyIn = buyIn
}

func (r *Round) Pot() uint {
	return r.pot
}

func (r *Round) SetPot(pot uint) {
	r.pot = pot
}

func (r *Round) Game() *LockStockGame {
	return r.game
}

func (r *Round) SetGame(game *LockStockGame) {
	r.game = game
}
