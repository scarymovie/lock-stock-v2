package model

type Bet struct {
	player *Player
	amount uint
	round  *Round
	number uint
}

func NewBet(player *Player, amount uint, round *Round, number uint) *Bet {
	return &Bet{player: player, amount: amount, round: round, number: number}
}
