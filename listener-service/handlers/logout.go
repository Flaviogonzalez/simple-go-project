package handlers

import (
	"encoding/json"
	"listener-service/helpers"
)

func LogoutHandler(data json.RawMessage) (r []byte, err error) {
	response, err := helpers.SendRequest("http://auth-service/logout", "POST", data)
	if err != nil {
		return nil, err
	}

	return response, nil
}
