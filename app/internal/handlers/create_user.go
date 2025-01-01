package handlers

import (
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
	w.WriteHeader(http.StatusOK)
	log.Printf("Player created %s", user.GetUserUid())
}
