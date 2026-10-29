// Package http — HTTP-эндпоинты Authorization Server OAuth 2.1 для MCP-клиентов: метаданные,
// динамическая регистрация, форма входа с выдачей кода, обмен кода и refresh-токена
package http

import (
	"context"
	stdErrors "errors"
	"net/http"

	"github.com/google/uuid"

	pkgErrors "pkg/errors"
	"pkg/log"
	"server/internal/modules/oauth/model"
	oauthService "server/internal/modules/oauth/service"
)

var _ OAuthService = new(oauthService.OAuthService)

// OAuthService — интерфейс Authorization Server
type OAuthService interface {
	PublicURL() string
	RegisterClient(context.Context, model.RegisterClientReq) (model.RegisterClientRes, error)
	ValidateAuthorize(context.Context, model.AuthorizeReq) (model.Client, error)
	GrantCode(context.Context, model.GrantCodeReq) (string, error)
	ExchangeCode(context.Context, model.ExchangeCodeReq) (model.TokenRes, error)
	Refresh(ctx context.Context, refreshToken string) (model.TokenRes, error)
}

// OAuthEndpoint — HTTP-обработчики OAuth
type OAuthEndpoint struct {
	oauthService OAuthService
}

// NewOAuthEndpoint создает обработчики OAuth
func NewOAuthEndpoint(oauthService OAuthService) *OAuthEndpoint {
	return &OAuthEndpoint{oauthService: oauthService}
}

// Register добавляет маршруты OAuth в mux
func (e *OAuthEndpoint) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /.well-known/oauth-authorization-server", e.Metadata)
	mux.HandleFunc("GET /.well-known/oauth-authorization-server/", e.Metadata)
	mux.HandleFunc("POST /oauth/register", e.RegisterClient)
	mux.HandleFunc("GET /oauth/authorize", e.AuthorizeForm)
	mux.HandleFunc("POST /oauth/authorize", e.AuthorizeSubmit)
	mux.HandleFunc("POST /oauth/token", e.Token)
}

// parseUUID разбирает client_id, пустая строка — uuid.Nil
func parseUUID(value string) uuid.UUID {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil
	}
	return id
}

// oauthError извлекает ошибку протокола; любая другая ошибка — server_error с записью в лог
func oauthError(err error) *model.OAuthError {
	var oauthErr *model.OAuthError
	if stdErrors.As(err, &oauthErr) {
		return oauthErr
	}
	log.LogError(pkgErrors.CastError(err))
	return model.NewOAuthError("server_error", "internal error")
}
