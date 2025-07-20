package handlers

import (
	"encoding/json"
	"fmt"
	"listener-service/helpers"
)

type ResponsePayload struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func HandleRegister(payload json.RawMessage) ([]byte, error) {
	response, err := helpers.SendRequest("http://auth-service/register", "POST", payload)
	if err != nil {
		return nil, err
	}

	var registerPayload ResponsePayload
	if err := json.Unmarshal(response, &registerPayload); err != nil {
		return nil, err
	}

	if registerPayload.Error {
		return nil, fmt.Errorf("registration failed: %s", registerPayload.Message)
	}

	return response, nil
}
