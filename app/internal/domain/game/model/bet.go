package model

type Bet struct {
	player *Player
	amount int
	round  *Round
}

func NewBet(player *Player, amount int, round *Round) *Bet {
	return &Bet{player: player, amount: amount, round: round}
}

func (b *Bet) Player() *Player {
	return b.player
}

func (b *Bet) SetPlayer(player *Player) {
	b.player = player
}

func (b *Bet) Amount() int {
	return b.amount
}

func (b *Bet) SetAmount(amount int) {
	b.amount = amount
}

func (b *Bet) Round() *Round {
	return b.round
}

func (b *Bet) SetRound(round *Round) {
	b.round = round
}
