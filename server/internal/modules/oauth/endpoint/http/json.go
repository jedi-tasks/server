package http

import (
	"encoding/json"
	"net/http"
)

// writeJSON пишет JSON-ответ с нужным статусом
func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// writeOAuthError пишет ошибку протокола (RFC 6749 §5.2)
func writeOAuthError(w http.ResponseWriter, err error) {
	oauthErr := oauthError(err)
	status := http.StatusBadRequest
	switch oauthErr.Code {
	case "invalid_client":
		status = http.StatusUnauthorized
	case "server_error":
		status = http.StatusInternalServerError
	}
	writeJSON(w, status, map[string]string{"error": oauthErr.Code, "error_description": oauthErr.Description})
}
