package domain

type RoomUserRepository interface {
	Save(roomUser RoomUser) error
	FindByRoom(room Room) ([]RoomUser, error)
}
