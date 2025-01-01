package domain

type Room struct {
	Id  int
	Uid string
}

func (r Room) GetRoomUid() string {
	return r.Uid
}
func (r Room) GetRoomId() int {
	return r.Id
}
