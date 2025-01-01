package domain

import "lock-stock-v2/external/domain"

type RoomUser struct {
	id   int
	room domain.Room
	user domain.User
}

func (ru *RoomUser) SetUser(user domain.User) {
	ru.user = user
}

func (ru *RoomUser) SetRoom(room domain.Room) {
	ru.room = room
}

func (ru *RoomUser) GetRoom() domain.Room {
	return ru.room
}

func (ru *RoomUser) GetUser() domain.User {
	return ru.user
}
