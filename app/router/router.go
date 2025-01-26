package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"lock-stock-v2/external/domain"
	"lock-stock-v2/external/handlers"
	"lock-stock-v2/middleware"
	"net/http"
)

func NewRouter(
	joinRoom handlers.JoinRoom,
	getAllRooms handlers.GetRooms,
	wsHandler handlers.WebSocketHandler,
	userFinder domain.UserFinder,
	createUser handlers.CreateUser,
	startGame handlers.StartGame,
) http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.With(
		middleware.LoggingMiddleware,
		middleware.UserAuthMiddleware(userFinder),
	).Route("/room", func(r chi.Router) {
		r.Post("/join/{roomId}", joinRoom.ServeHTTP)
		r.Post("/list", getAllRooms.ServeHTTP)
		r.Post("/start/{roomId}", startGame.ServeHTTP)
	})

	r.With(
		middleware.LoggingMiddleware,
	).Post("/user/create", createUser.ServeHTTP)

	r.Handle("/ws/{roomId}", wsHandler)

	return r
}
