package usecase

type JoinRoom interface {
	JoinRoom(request JoinRoomRequest) error
}
