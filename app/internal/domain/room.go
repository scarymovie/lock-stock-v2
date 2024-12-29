package domain

type Room struct {
	Id string
}

func (r Room) GetId() string {
	return r.Id
}
