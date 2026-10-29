package http

import "net/http"

// Metadata отдает метаданные Authorization Server (RFC 8414): клиенты находят по ним эндпоинты
func (e *OAuthEndpoint) Metadata(w http.ResponseWriter, _ *http.Request) {
	issuer := e.oauthService.PublicURL()
	writeJSON(w, http.StatusOK, map[string]any{
		"issuer":                                issuer,
		"authorization_endpoint":                issuer + "/oauth/authorize",
		"token_endpoint":                        issuer + "/oauth/token",
		"registration_endpoint":                 issuer + "/oauth/register",
		"response_types_supported":              []string{"code"},
		"response_modes_supported":              []string{"query"},
		"grant_types_supported":                 []string{"authorization_code", "refresh_token"},
		"code_challenge_methods_supported":      []string{"S256"},
		"token_endpoint_auth_methods_supported": []string{"none"},
		"scopes_supported":                      []string{},
	})
}
