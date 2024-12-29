package domain

type RoomRepository interface {
	Save(room Room) error
}
