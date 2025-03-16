package model

type RoundPlayerLog struct {
	player    *Player
	round     *Round
	number    uint
	betsValue uint
	answer    *uint
}

func NewRoundPlayerLog(player *Player, round *Round, number uint, betsValue uint) *RoundPlayerLog {
	return &RoundPlayerLog{player: player, round: round, number: number, betsValue: betsValue}
}

func (r *RoundPlayerLog) Answer() *uint {
	return r.answer
}

func (r *RoundPlayerLog) SetAnswer(answer *uint) {
	r.answer = answer
}

func (r *RoundPlayerLog) SetBetsValue(betsValue uint) {
	r.betsValue = betsValue
}

func (r *RoundPlayerLog) Player() *Player {
	return r.player
}

func (r *RoundPlayerLog) Round() *Round {
	return r.round
}

func (r *RoundPlayerLog) Number() uint {
	return r.number
}

func (r *RoundPlayerLog) BetsValue() uint {
	return r.betsValue
}
