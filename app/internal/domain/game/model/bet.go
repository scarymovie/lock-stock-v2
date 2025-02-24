package model

type Bet struct {
	player *Player
	amount uint
	round  *Round
}
