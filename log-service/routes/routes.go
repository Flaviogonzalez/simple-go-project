package routes

import (
	"database/sql"
	"log-service/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes(db *sql.DB) http.Handler {
	mux := chi.NewRouter()

	// en mi cabeza se veia mucho mejor voy a tener que retomar en un tiempo

	mux.Post("/log/info", handlers.InfoLogHandler(db))
	mux.Post("/log/debug", handlers.DebugLogHandler(db))
	mux.Post("/log/warning", handlers.WarningLogHandler(db))
	mux.Post("/log/error", handlers.ErrorLogHandler(db))
	mux.Post("/log/fatal", handlers.FatalLogHandler(db))
	return mux
}
