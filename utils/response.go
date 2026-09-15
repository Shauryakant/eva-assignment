package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func JSONError(w http.ResponseWriter, statusCode int, message string) {
	JSON(w, statusCode, ErrorResponse{Error: message})
}
