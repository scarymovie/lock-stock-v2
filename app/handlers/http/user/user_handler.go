package user

import (
	"encoding/json"
	"fmt"
	"lock-stock-v2/internal/domain/user/service"
	"log"
	"net/http"
)

type UserHandler struct {
	createUser *service.CreateUserService
}

func NewUserHandler(createUser *service.CreateUserService) *UserHandler {
	return &UserHandler{createUser: createUser}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var rawUser service.RawCreateUser

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

	var userUid = user.Uid()
	var userName = user.Name()

	var responseData = User{
		UserId: &userUid,
		Name:   &userName,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(responseData); err != nil {
		log.Printf("Ошибка при отправке ответа: %v", err)
	}

	log.Printf("Player created %s", user.Uid())
}
