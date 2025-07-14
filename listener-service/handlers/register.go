package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ResponsePayload struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func HandleRegister(payload json.RawMessage) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://auth-service/register", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	var registerPayload ResponsePayload
	err = json.Unmarshal(response, &registerPayload)
	if err != nil {
		return response, err
	}

	if registerPayload.Error {
		return response, fmt.Errorf("registration failed: %s", registerPayload.Message)
	}

	return response, nil
}
