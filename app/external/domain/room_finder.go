package domain

type RoomFinder interface {
	FindById(roomId string) (Room, error)
	GetAll() ([]Room, error)
}
