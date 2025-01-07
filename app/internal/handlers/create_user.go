package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"lock-stock-v2/external/domain"
	internalDomain "lock-stock-v2/internal/domain"
	"log"
	"net/http"
)

type CreateUser struct {
	userRepository domain.UserRepository
}

func NewCreateUser(userRepository domain.UserRepository) *CreateUser {
	return &CreateUser{userRepository: userRepository}
}

func (h *CreateUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var user = internalDomain.User{
		Uid: uuid.New().String(),
	}
	err := h.userRepository.SaveUser(user)
	if err != nil {
		return
	}
	response := map[string]string{
		"user_id": user.GetUserUid(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}

	log.Printf("Player created %s", user.GetUserUid())
}
