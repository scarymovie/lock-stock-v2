package domain

type Room struct {
	Id string
}

func (r Room) GetRoomId() string {
	return r.Id
}
