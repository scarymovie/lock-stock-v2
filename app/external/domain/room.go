package domain

type Room interface {
	Uid() string
	Status() RoomStatus
	SetStatus(roomStatus RoomStatus)
}
