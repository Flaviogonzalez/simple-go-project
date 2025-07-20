package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type LogInfo struct {
	Message           string `json:"message"`
	Thread_identifier string `json:"thread_identifier"`
	Requestid         string `json:"requestid"`
	Userid            string `json:"userid"`
}

func InfoLogHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var log LogInfo

		err := json.NewDecoder(r.Body).Decode(&log)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

	}
}
