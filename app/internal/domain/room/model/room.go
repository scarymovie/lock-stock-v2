package model

type Room struct {
	uid    string
	status RoomStatus
}

func NewRoom(uid string, status RoomStatus) *Room {
	return &Room{uid: uid, status: status}
}

func (r *Room) Uid() string {
	return r.uid
}

func (r *Room) Status() RoomStatus {
	return r.status
}

func (r *Room) SetStatus(status RoomStatus) {
	r.status = status
}
