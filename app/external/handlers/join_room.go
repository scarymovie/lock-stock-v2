package handlers

import "net/http"

type JoinRoom interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
