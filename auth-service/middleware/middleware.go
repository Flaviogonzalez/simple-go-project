package middleware

import (
	"context"
	"database/sql"
	"net/http"
)

type CtxKey string

const DbKey CtxKey = "db"

func HandlerWrapper(db *sql.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), DbKey, db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// To retrieve db from context in other routes:
// db, ok := r.Context().Value(dbKey).(*sql.DB)
