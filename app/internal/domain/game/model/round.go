package model

type Round struct {
	number   *uint
	question *Question
	buyIn    uint
	pot      uint
	game     *LockStockGame
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

func NewRound(number *uint, question *Question, buyIn uint, game *LockStockGame) *Round {
	return &Round{number: number, question: question, buyIn: buyIn, game: game}
}
