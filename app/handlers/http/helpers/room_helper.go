package helpers

import (
	"lock-stock-v2/internal/domain/room/model"
	"lock-stock-v2/internal/domain/room/repository"
	"log"
	"net/http"
)

func GetRoomById(roomRepository repository.RoomRepository, roomID string) (*model.Room, error) {
	if roomID == "" {
		return nil, &RoomNotFoundError{Code: http.StatusBadRequest, Message: "Room ID is required"}
	}

	room, err := roomRepository.FindById(roomID)
	if err != nil {
		log.Println("FindById error:", err)
		return nil, &RoomNotFoundError{Code: http.StatusNotFound, Message: "Failed to find room"}
	}

	if room.Uid() == "" {
		log.Printf("Room %v does not exist or invalid\n", roomID)
		return nil, &RoomNotFoundError{Code: http.StatusNotFound, Message: "Room does not exist or invalid"}
	}

	log.Printf("Room retrieved: UID=%s", room.Uid())
	return &room, nil
}

func (e *RoomNotFoundError) Error() string {
	return e.Message
}

type RoomNotFoundError struct {
	Code    int
	Message string
}
