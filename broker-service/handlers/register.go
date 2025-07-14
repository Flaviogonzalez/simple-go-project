package handlers

import (
	"broker-service/event"
	"encoding/json"
	"io"
	"net/http"
)

type RegisterPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Policy   int32  `json:"policy"`
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var registerPayload RegisterPayload
	err = json.Unmarshal(body, &registerPayload)
	if err != nil {
		http.Error(w, "Failed to unmarshal payload", http.StatusBadRequest)
		return
	}

	if registerPayload.Email == "" || registerPayload.Password == "" || registerPayload.Name == "" {
		http.Error(w, "Name, email, and password are required", http.StatusBadRequest)
		return
	}

	if registerPayload.Policy == 0 {
		http.Error(w, "Must accept the privacy policy", http.StatusBadRequest)
		return
	}

	payloadBytes, err := json.Marshal(registerPayload)
	if err != nil {
		http.Error(w, "Failed to marshal payload", http.StatusInternalServerError)
		return
	}

	topicPayload := event.TopicPayload{
		Name: "user.register",
		Event: event.EventPayload{
			Name: "RegisterEvent",
			Data: json.RawMessage(payloadBytes),
		},
	}

	err = event.SendToListener(w, "AuthenticationService", topicPayload)
	if err != nil {
		http.Error(w, "Failed to send registration event: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
