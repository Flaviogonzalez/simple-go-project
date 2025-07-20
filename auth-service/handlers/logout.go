package handlers

import (
	"auth-service/helpers"
	"auth-service/models"
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
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

	err = q.DeleteSession(r.Context(), sessionID)
	if err != nil {
		helpers.ErrorJSON(w, http.StatusInternalServerError, "error deleting session")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now().Add(-time.Second * 1),
		HttpOnly: true,
		Secure:   true,
	})

	response := helpers.ErrorPayload{
		Error:   false,
		Message: "Logout successful",
	}

	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")

	helpers.WriteJSON(w, http.StatusOK, response, headers)
}
