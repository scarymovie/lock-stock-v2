package model

type Round struct {
	number   *uint
	question *Question
	buyIn    uint
	game     *LockStockGame
}

func NewRound(number *uint, question *Question, buyIn uint, game *LockStockGame) *Round {
	return &Round{number: number, question: question, buyIn: buyIn, game: game}
}
