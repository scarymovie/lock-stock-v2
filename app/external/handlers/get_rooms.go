package handlers

import "net/http"

type GetRooms interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
