package http

import (
	"encoding/json"
	"net/http"

	"server/internal/modules/oauth/model"
)

// RegisterClient регистрирует клиента динамически (RFC 7591)
func (e *OAuthEndpoint) RegisterClient(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterClientReq
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10)).Decode(&req); err != nil {
		writeOAuthError(w, model.NewOAuthError("invalid_client_metadata", "body must be JSON"))
		return
	}

	res, err := e.oauthService.RegisterClient(r.Context(), req)
	if err != nil {
		writeOAuthError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, res)
}
