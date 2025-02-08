package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"lock-stock-v2/handlers/http/room"
	"lock-stock-v2/handlers/http/user"
	"lock-stock-v2/handlers/http/ws"
	"lock-stock-v2/internal/domain/user/repository"
	"lock-stock-v2/middleware"
	"log"
	"net/http"
)

func NewRouter(
	roomHandler room.RoomHandler,
	userHandler user.UserHandler,
	wsHandler ws.WebSocketHandler,
	userRepository repository.UserRepository,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.LoggingAllRequests)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Group(func(r chi.Router) {
		r.With(
			middleware.LoggingMiddleware,
			middleware.UserAuthMiddleware(userRepository),
		).Route("/room", func(r chi.Router) {
			room.HandlerFromMux(&roomHandler, r)
		})
	})

	r.Group(func(r chi.Router) {
		r.With(middleware.LoggingMiddleware).Route("/user", func(r chi.Router) {
			user.HandlerFromMux(&userHandler, r)
		})
	})

	r.Group(func(r chi.Router) {
		r.With(middleware.LoggingMiddleware).Route("/ws/{roomId}", func(r chi.Router) {
			ws.HandlerFromMux(&wsHandler, r)
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Not found url %s\n", r.URL)
		w.WriteHeader(http.StatusNotFound)
	})

	return r
}
