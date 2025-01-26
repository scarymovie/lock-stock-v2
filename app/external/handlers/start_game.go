package handlers

import "net/http"

type StartGame interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
