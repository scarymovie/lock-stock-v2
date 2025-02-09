package model

import (
	roomModel "lock-stock-v2/internal/domain/room/model"
	userModel "lock-stock-v2/internal/domain/user/model"
)

type RoomUser struct {
	room *roomModel.Room
	user *userModel.User
}

func NewRoomUser(room *roomModel.Room, user *userModel.User) *RoomUser {
	return &RoomUser{room: room, user: user}
}

func (ru *RoomUser) SetUser(user *userModel.User) {
	ru.user = user
}

func (ru *RoomUser) SetRoom(room *roomModel.Room) {
	ru.room = room
}

func (ru *RoomUser) Room() *roomModel.Room {
	return ru.room
}

func (ru *RoomUser) User() *userModel.User {
	return ru.user
}
