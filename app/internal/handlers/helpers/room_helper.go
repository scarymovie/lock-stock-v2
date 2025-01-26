package helpers

import (
	"lock-stock-v2/external/domain"
	"log"
	"net/http"
)

func GetRoomById(roomFinder domain.RoomFinder, roomID string) (domain.Room, error) {
	if roomID == "" {
		return nil, &RoomNotFoundError{Code: http.StatusBadRequest, Message: "Room ID is required"}
	}

	room, err := roomFinder.FindById(roomID)
	if err != nil {
		log.Println("FindById error:", err)
		return nil, &RoomNotFoundError{Code: http.StatusNotFound, Message: "Failed to find room"}
	}

	if room == nil || room.GetRoomId() == 0 {
		log.Printf("Room %v does not exist or invalid\n", roomID)
		return nil, &RoomNotFoundError{Code: http.StatusNotFound, Message: "Room does not exist or invalid"}
	}

	log.Printf("Room retrieved: ID=%d, UID=%s", room.GetRoomId(), room.GetRoomUid())
	return room, nil
}

func (e *RoomNotFoundError) Error() string {
	return e.Message
}

type RoomNotFoundError struct {
	Code    int
	Message string
}
