package handlers

import (
	"broker-service/event"
	"encoding/json"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	SessionID string `json:"session_id"`
	Token     string `json:"token"`
	CsrfToken string `json:"csrf_token"`
	Error     bool   `json:"error"`
	Message   string `json:"message"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Failed to unmarshal payload", http.StatusBadRequest)
		return
	}

	if request.Email == "" || request.Password == "" {
		http.Error(w, "asdsasdas", http.StatusBadRequest)
		return
	}

	payload, err := json.Marshal(request)
	if err != nil {
		http.Error(w, "Failed to marshal payload", http.StatusBadRequest)
		return
	}

	topicPayload := event.TopicPayload{
		Name: "user.login",
		Event: event.EventPayload{
			Name: "LoginEvent",
			Data: payload,
		},
	}

	err = event.SendToListener(w, "AuthenticationService", topicPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
