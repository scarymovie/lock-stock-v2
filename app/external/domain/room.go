package domain

type Room interface {
	GetRoomUid() string
	GetRoomId() int
	GetRoomStatus() RoomStatus
	SetRoomStatus(roomStatus RoomStatus)
}
