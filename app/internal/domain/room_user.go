package domain

import "lock-stock-v2/external/domain"

type RoomUser struct {
	room domain.Room
	user domain.User
}

func NewRoomUser(room domain.Room, user domain.User) *RoomUser {
	return &RoomUser{room: room, user: user}
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
