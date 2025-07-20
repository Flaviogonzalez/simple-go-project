package middleware

import (
	"auth-service/helpers"
	"auth-service/models"
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func CSRFCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session_id")
		if err != nil {
			helpers.ErrorJSON(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		db := r.Context().Value("db").(*sql.DB)
		q := models.New(db)

		sessionID, err := uuid.Parse(cookie.Value)
		if err != nil {
			helpers.ErrorJSON(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		session, err := q.GetSessionByID(r.Context(), sessionID)
		if err != nil {
			helpers.ErrorJSON(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		if session.Expires.Before(time.Now()) {
			helpers.ErrorJSON(w, http.StatusUnauthorized, "invalid or expired session")
			return
		}

		csrfToken := r.Header.Get("X-CSRF-Token")
		if csrfToken == "" {
			helpers.ErrorJSON(w, http.StatusForbidden, "csrf token is required")
			return
		}

		if session.CsrfToken != csrfToken {
			helpers.ErrorJSON(w, http.StatusForbidden, "invalid csrf token")
			return
		}

		next.ServeHTTP(w, r)
	})
}
