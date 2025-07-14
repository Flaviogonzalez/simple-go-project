package helpers

import (
	"encoding/json"
	"net/http"
)

type ErrorPayload struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	return json.NewDecoder(r.Body).Decode(data)
}

func ErrorJSON(w http.ResponseWriter, errorCode int, message string) {
	if errorCode == 0 {
		errorCode = http.StatusInternalServerError
	}

	payload := ErrorPayload{
		Error:   true,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		// Fallback to a generic error if marshaling fails
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": true, "message": "Internal server errorr"}`))
		return
	}

	w.WriteHeader(errorCode)
	w.Write(data)
}
