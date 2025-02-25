package model

type Bet struct {
	player *Player
	amount uint
	round  *Round
}

func NewBet(player *Player, amount uint, round *Round) *Bet {
	return &Bet{player: player, amount: amount, round: round}
}
