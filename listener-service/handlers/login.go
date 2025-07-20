package handlers

import (
	"encoding/json"
	"fmt"
	"listener-service/helpers"
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

func LoginHandler(data json.RawMessage) (r []byte, err error) {
	var request LoginRequest

	if err := json.Unmarshal(data, &request); err != nil {
		return nil, err
	}

	response, err := helpers.SendRequest("http://auth-service/login", "POST", data)
	if err != nil {
		return nil, err
	}

	var loginResponse LoginResponse
	if err := json.Unmarshal(response, &loginResponse); err != nil {
		return nil, err
	}

	if loginResponse.Error {
		return nil, fmt.Errorf("login failed: %s", loginResponse.Message)
	}

	return response, nil
}
