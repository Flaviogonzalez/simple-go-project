package routes

import (
	"broker-service/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://*", "https://*"},
		AllowedMethods: []string{"POST", "PUT", "DELETE", "GET", "OPTIONS"},
		ExposedHeaders: []string{"Content-Type", "Authorization", "X-CSRF-Token"},
		MaxAge:         300,
	}))

	mux.Post("/register", handlers.HandleRegister)

	return mux
}
