package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ty16akin/ConfirmNG/internal/handlers"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/users", loadUserRoutes)

	return router
}

func loadUserRoutes(router chi.Router) {
	userHandler := &handlers.User{}
	router.Post("/", userHandler.CreateUser)
	router.Get("/", userHandler.ListUser)
	router.Get("/{id}", userHandler.GetUserByID)
	router.Put("/{id}", userHandler.UpdateUserbyID)
	router.Delete("/{id}", userHandler.DeleteUserById)
}
