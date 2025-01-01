package domain

type RoomUser interface {
	GetRoom() Room
	GetUser() User
	SetRoom(Room)
	SetUser(User)
}
