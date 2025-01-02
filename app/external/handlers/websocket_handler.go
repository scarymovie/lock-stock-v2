package handlers

import "net/http"

type WebSocketHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
