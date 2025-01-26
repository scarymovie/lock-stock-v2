package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	externalDomain "lock-stock-v2/external/domain"
)

type GetRooms struct {
	roomFinder externalDomain.RoomFinder
}

func NewGetRooms(roomFinder externalDomain.RoomFinder) *GetRooms {
	return &GetRooms{roomFinder: roomFinder}
}

func (g *GetRooms) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rooms, err := g.roomFinder.GetPending()
	if err != nil {
		http.Error(w, "Failed to get rooms: "+err.Error(), http.StatusInternalServerError)
		return
	}

	type roomResponse struct {
		RoomUid string `json:"roomUid"`
	}

	var responseData []roomResponse
	for _, room := range rooms {
		responseData = append(responseData, roomResponse{
			RoomUid: room.GetRoomUid(), // замените на GetUID(), если у вас другой метод
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		log.Printf("Ошибка при кодировании ответа: %v\n", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
