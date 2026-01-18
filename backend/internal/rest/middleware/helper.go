package middleware

import (
	"encoding/json"
	"net/http"
)


type ErrorResponse struct {
    Error string `json:"error"`
}

func sendJSON(w http.ResponseWriter, status int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(payload)
}

func sendError(w http.ResponseWriter, status int, message string) {
    sendJSON(w, status, ErrorResponse{Error: message})
}
