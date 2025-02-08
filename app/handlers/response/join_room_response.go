package response

type RoomUserResponse struct {
	RoomUid  string `json:"room_uid"`
	UserUid  string `json:"user_uid"`
	UserName string `json:"user_name"`
}
