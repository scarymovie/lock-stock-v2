package usecase

type JoinRoomRequest struct {
	PlayerId string `json:"playerId"`
	RoomId   string `json:"roomId"`
}
