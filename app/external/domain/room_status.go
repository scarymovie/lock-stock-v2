package domain

type RoomStatus string

const (
	StatusPending  RoomStatus = "pending"
	StatusStarted  RoomStatus = "started"
	StatusFinished RoomStatus = "finished"
	StatusCanceled RoomStatus = "canceled"
)
