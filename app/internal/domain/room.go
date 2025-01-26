package domain

import "lock-stock-v2/external/domain"

type Room struct {
	Id     int
	Uid    string
	Status domain.RoomStatus
}

func (r *Room) GetRoomUid() string {
	return r.Uid
}
func (r *Room) GetRoomId() int {
	return r.Id
}

func (r *Room) GetRoomStatus() domain.RoomStatus {
	return r.Status
}

func (r *Room) SetRoomStatus(status domain.RoomStatus) {
	r.Status = status
}
