package handlers

import "net/http"

type CreateUser interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
