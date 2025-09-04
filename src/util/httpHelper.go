package util

import (
	"contatos/model"
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type WebHandlerFunc func(db *sqlx.DB, w http.ResponseWriter, r *http.Request, c Config, userID string)

func RespondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := model.APIResponse{
		Success: false,
		Error:   message,
	}

	json.NewEncoder(w).Encode(response)
}

func RespondWithSuccess(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := model.APIResponse{
		Success: true,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}
