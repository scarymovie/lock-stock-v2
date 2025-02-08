package handlers

import (
	"encoding/json"
	"fmt"
	"lock-stock-v2/external/usecase"
	"log"
	"net/http"
)

type CreateUser struct {
	createUser usecase.CreateUser
}

func NewCreateUser(createUser usecase.CreateUser) *CreateUser {
	return &CreateUser{createUser: createUser}
}

func (h *CreateUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var rawUser usecase.RawCreateUser

	if err := json.NewDecoder(r.Body).Decode(&rawUser); err != nil {
		http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.createUser.Do(rawUser)
	if err != nil {
		fmt.Printf("Ошибка при создании пользователя: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"user_id": user.Uid(),
		"name":    user.Name(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}

	log.Printf("Player created %s", user.Uid())
}
