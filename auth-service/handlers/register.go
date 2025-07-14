package handlers

import (
	"auth-service/helpers"
	"auth-service/middleware"
	"auth-service/models"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RegisterPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Policy   int32  `json:"policy"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var registerPayload RegisterPayload
	if err := helpers.ReadJSON(w, r, &registerPayload); err != nil {
		helpers.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if registerPayload.Email == "" || registerPayload.Password == "" || registerPayload.Name == "" {
		helpers.ErrorJSON(w, http.StatusBadRequest, "Name, email, and password are required")
		return
	}

	if registerPayload.Policy == 0 {
		helpers.ErrorJSON(w, http.StatusBadRequest, "Must accept the privacy policy")
		return
	}

	db, ok := r.Context().Value(middleware.DbKey).(*sql.DB)
	if !ok {
		helpers.ErrorJSON(w, http.StatusInternalServerError, "Database connection error")
		return
	}

	q := models.New(db)

	user, err := q.GetUserByEmail(r.Context(), registerPayload.Email)
	if err != nil && err != sql.ErrNoRows {
		helpers.ErrorJSON(w, http.StatusInternalServerError, "Error checking email existence")
		return
	}

	if user.Email == registerPayload.Email {
		helpers.ErrorJSON(w, http.StatusConflict, "Email already registered")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerPayload.Password), bcrypt.DefaultCost)
	if err != nil {
		helpers.ErrorJSON(w, http.StatusInternalServerError, "Error processing password")
		return
	}

	_, err = q.CreateUser(r.Context(), models.CreateUserParams{
		ID:       uuid.New(),
		Name:     registerPayload.Name,
		Email:    registerPayload.Email,
		Password: string(hashedPassword),
		Policy:   registerPayload.Policy,
	})
	if err != nil {
		helpers.ErrorJSON(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	response, err := json.Marshal(helpers.ErrorPayload{
		Error:   false,
		Message: "User registered successfully",
	})
	if err != nil {
		helpers.ErrorJSON(w, http.StatusInternalServerError, "Error creating response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
