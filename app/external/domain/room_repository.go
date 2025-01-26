package domain

type RoomRepository interface {
	Save(room Room) error
	UpdateRoomStatus(room Room) error
}
