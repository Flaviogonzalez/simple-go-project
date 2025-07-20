package handlers

import (
	"auth-service/helpers"
	"auth-service/middleware"
	"auth-service/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RequestPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResponsePayload struct {
	SessionID string `json:"session_id"`
	Token     string `json:"token"`      // JWT token
	CsrfToken string `json:"csrf_token"` // CSRF token
	Error     bool   `json:"error"`
	Message   string `json:"message"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var request RequestPayload
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		helpers.ErrorJSON(w, http.StatusBadRequest, "error reading request")
		return
	}

	if request.Email == "" || request.Password == "" {
		helpers.ErrorJSON(w, http.StatusUnauthorized, "Inputs must contain Email and Password")
		return
	}

	db, ok := r.Context().Value(middleware.DbKey).(*sql.DB)
	if !ok {
		helpers.ErrorJSON(w, http.StatusInternalServerError, "internal server error")
		return
	}

	q := models.New(db)

	user, err := q.GetUserByEmail(r.Context(), request.Email)
	if err != nil {
		helpers.ErrorJSON(w, http.StatusUnauthorized, "error getting further information")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		helpers.ErrorJSON(w, http.StatusUnauthorized, err.Error())
		return
	}

	// create session
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		helpers.ErrorJSON(w, http.StatusInternalServerError, "error signing token")
		return
	}

	session, err := q.CreateSession(r.Context(), models.CreateSessionParams{
		ID:           uuid.New(),
		UserID:       user.ID,
		SessionToken: uuid.New().String(),
		CsrfToken:    uuid.New().String(),
		Expires:      time.Now().Add(time.Hour * 24),
	})
	if err != nil {
		helpers.ErrorJSON(w, http.StatusInternalServerError, "error creating session")
		return
	}
	// store session in cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.SessionToken,
		Expires:  session.Expires,
		HttpOnly: true,
		Secure:   true,
	})

	// generate
	response := ResponsePayload{
		SessionID: session.ID.String(),
		Token:     tokenString,
		CsrfToken: session.CsrfToken,
		Error:     false,
		Message:   "Login successful",
	}

	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")

	helpers.WriteJSON(w, http.StatusOK, response, headers)
}
