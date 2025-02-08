package domain

import "lock-stock-v2/external/domain"

type Room struct {
	uid    string
	status domain.RoomStatus
}

func NewRoom(uid string, status domain.RoomStatus) *Room {
	return &Room{uid: uid, status: status}
}

func (r *Room) Uid() string {
	return r.uid
}

func (r *Room) Status() domain.RoomStatus {
	return r.status
}

func (r *Room) SetStatus(status domain.RoomStatus) {
	r.status = status
}
