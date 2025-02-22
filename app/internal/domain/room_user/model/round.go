package model

type Round struct {
	roomUser *RoomUser
	number   int
}

func NewRound(roomUser *RoomUser, number int) *Round {
	return &Round{roomUser: roomUser, number: number}
}
