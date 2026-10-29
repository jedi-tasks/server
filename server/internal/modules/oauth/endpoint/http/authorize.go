package http

import (
	stdErrors "errors"
	"net/http"
	"net/url"

	"server/internal/modules/oauth/model"
)

// authorizeReq собирает параметры авторизации из query или полей формы
func authorizeReq(values url.Values) model.AuthorizeReq {
	return model.AuthorizeReq{
		ClientID:            parseUUID(values.Get("client_id")),
		RedirectURI:         values.Get("redirect_uri"),
		ResponseType:        values.Get("response_type"),
		CodeChallenge:       values.Get("code_challenge"),
		CodeChallengeMethod: values.Get("code_challenge_method"),
		Resource:            values.Get("resource"),
		State:               values.Get("state"),
	}
}

// AuthorizeForm проверяет запрос авторизации и показывает форму входа
func (e *OAuthEndpoint) AuthorizeForm(w http.ResponseWriter, r *http.Request) {
	req := authorizeReq(r.URL.Query())
	client, err := e.oauthService.ValidateAuthorize(r.Context(), req)
	if err != nil {
		e.authorizeError(w, req, err)
		return
	}
	renderForm(w, http.StatusOK, formData{Client: client.Name, Params: r.URL.Query(), Error: ""})
}

// AuthorizeSubmit проверяет пароль и редиректит клиенту с кодом; при неверном пароле — форма снова
func (e *OAuthEndpoint) AuthorizeSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeOAuthError(w, model.NewOAuthError("invalid_request", "bad form"))
		return
	}
	req := authorizeReq(r.PostForm)

	code, err := e.oauthService.GrantCode(r.Context(), model.GrantCodeReq{
		Authorize: req,
		Email:     r.PostForm.Get("email"),
		Password:  r.PostForm.Get("password"),
	})
	if err != nil {
		// Неверный пароль — остаемся на форме, остальное — как при GET
		if stdErrors.Is(err, model.ErrInvalidCredentials) {
			client, validateErr := e.oauthService.ValidateAuthorize(r.Context(), req)
			if validateErr != nil {
				e.authorizeError(w, req, validateErr)
				return
			}
			renderForm(w, http.StatusUnauthorized, formData{Client: client.Name, Params: r.PostForm, Error: "Неверная почта или пароль"})
			return
		}
		e.authorizeError(w, req, err)
		return
	}

	// Успех: код и state — клиенту
	redirect, _ := url.Parse(req.RedirectURI)
	query := redirect.Query()
	query.Set("code", code)
	if req.State != "" {
		query.Set("state", req.State)
	}
	redirect.RawQuery = query.Encode()
	http.Redirect(w, r, redirect.String(), http.StatusFound)
}

// authorizeError: ошибку клиента или redirect_uri показываем пользователю, остальные — редиректом
// клиенту (RFC 6749 §4.1.2.1)
func (e *OAuthEndpoint) authorizeError(w http.ResponseWriter, req model.AuthorizeReq, err error) {
	oauthErr := oauthError(err)
	if oauthErr.Code == "invalid_client" || oauthErr.Code == "invalid_request" && req.RedirectURI == "" ||
		oauthErr.Description == "redirect_uri is not registered for this client" {
		renderError(w, http.StatusBadRequest, oauthErr.Description)
		return
	}
	redirect, parseErr := url.Parse(req.RedirectURI)
	if parseErr != nil {
		renderError(w, http.StatusBadRequest, oauthErr.Description)
		return
	}
	query := redirect.Query()
	query.Set("error", oauthErr.Code)
	query.Set("error_description", oauthErr.Description)
	if req.State != "" {
		query.Set("state", req.State)
	}
	redirect.RawQuery = query.Encode()
	w.Header().Set("Location", redirect.String())
	w.WriteHeader(http.StatusFound)
}
