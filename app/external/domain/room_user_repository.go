package domain

type RoomUserRepository interface {
	Save(roomUser RoomUser) error
}
