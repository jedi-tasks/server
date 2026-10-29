package http

import (
	"net/http"

	"server/internal/modules/oauth/model"
)

// Token обменивает код или refresh-токен на пару токенов (RFC 6749 §4.1.3, §6)
func (e *OAuthEndpoint) Token(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeOAuthError(w, model.NewOAuthError("invalid_request", "body must be application/x-www-form-urlencoded"))
		return
	}
	form := r.PostForm

	var res model.TokenRes
	var err error
	switch form.Get("grant_type") {
	case "authorization_code":
		res, err = e.oauthService.ExchangeCode(r.Context(), model.ExchangeCodeReq{
			Code:         form.Get("code"),
			ClientID:     form.Get("client_id"),
			RedirectURI:  form.Get("redirect_uri"),
			CodeVerifier: form.Get("code_verifier"),
			Resource:     form.Get("resource"),
		})
	case "refresh_token":
		res, err = e.oauthService.Refresh(r.Context(), form.Get("refresh_token"))
	default:
		err = model.NewOAuthError("unsupported_grant_type", "grant_type must be authorization_code or refresh_token")
	}
	if err != nil {
		writeOAuthError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}
