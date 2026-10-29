package service

import (
	"context"
	"slices"

	"server/internal/modules/oauth/model"
)

// ValidateAuthorize проверяет параметры запроса авторизации и возвращает клиента для формы входа.
// Ошибка клиента или redirect_uri — показывается пользователю; остальные ошибки endpoint
// возвращает редиректом на redirect_uri
func (s *OAuthService) ValidateAuthorize(ctx context.Context, req model.AuthorizeReq) (model.Client, error) {
	ctx, span := tracer.Start(ctx, "ValidateAuthorize")
	defer span.End()

	// Клиент и redirect_uri — без них редиректить некуда
	client, found, err := s.oauthRepository.GetClient(ctx, req.ClientID)
	if err != nil {
		return client, err
	}
	if !found {
		return client, model.NewOAuthError("invalid_client", "unknown client_id")
	}
	if !slices.Contains(client.RedirectURIs, req.RedirectURI) {
		return client, model.NewOAuthError("invalid_request", "redirect_uri is not registered for this client")
	}

	// Остальные параметры
	if req.ResponseType != "code" {
		return client, model.NewOAuthError("unsupported_response_type", "only response_type=code is supported")
	}
	if req.CodeChallenge == "" || req.CodeChallengeMethod != "S256" {
		return client, model.NewOAuthError("invalid_request", "PKCE with code_challenge_method=S256 is required")
	}
	if req.Resource != "" && !sameResource(req.Resource, s.config.PublicURL) {
		return client, model.NewOAuthError("invalid_target", "resource must be "+s.config.PublicURL)
	}

	return client, nil
}
