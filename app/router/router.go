package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"lock-stock-v2/external/domain"
	"lock-stock-v2/handlers"
	"lock-stock-v2/middleware"
	"log"
	"net/http"
)

func NewRouter(
	joinRoom *handlers.JoinRoom,
	getAllRooms *handlers.GetRooms,
	wsHandler *handlers.WebSocketHandler,
	createUser *handlers.CreateUser,
	startGame *handlers.StartGame,
	userFinder domain.UserFinder,
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

	r.With(
		middleware.LoggingMiddleware,
		middleware.UserAuthMiddleware(userFinder),
	).Route("/room", func(r chi.Router) {
		r.Post("/join/{roomId}", joinRoom.Do)
		r.Post("/list", getAllRooms.Do)
		r.Post("/start/{roomId}", startGame.Do)
	})

	r.With(
		middleware.LoggingMiddleware,
	).Post("/user/create", createUser.Do)

	r.Handle("/ws/{roomId}", wsHandler)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Not found url %s\n", r.URL)
		w.WriteHeader(http.StatusNotFound)
	})

	return r
}
