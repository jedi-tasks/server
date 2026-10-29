package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/google/uuid"

	"server/internal/enum/osType"
	authModel "server/internal/modules/auth/model"
	"server/internal/modules/oauth/model"
	userModel "server/internal/modules/user/model"
	"server/internal/utils/errors"
)

// ExchangeCode обменивает код на токены сервера: проверяет клиента, redirect_uri, PKCE и resource,
// создает для выдачи отдельное устройство без локальной копии данных
func (s *OAuthService) ExchangeCode(ctx context.Context, req model.ExchangeCodeReq) (model.TokenRes, error) {
	ctx, span := tracer.Start(ctx, "ExchangeCode")
	defer span.End()

	// Код одноразовый: ConsumeCode помечает его использованным атомарно
	code, found, err := s.oauthRepository.ConsumeCode(ctx, req.Code, time.Now())
	if err != nil {
		return model.TokenRes{}, err
	}
	if !found {
		return model.TokenRes{}, model.NewOAuthError("invalid_grant", "code is unknown, expired or already used")
	}

	// Код должен обменивать тот же клиент на тот же адрес
	clientID, err := uuid.Parse(req.ClientID)
	if err != nil || clientID != code.ClientID {
		return model.TokenRes{}, model.NewOAuthError("invalid_grant", "client_id does not match the code")
	}
	if req.RedirectURI != code.RedirectURI {
		return model.TokenRes{}, model.NewOAuthError("invalid_grant", "redirect_uri does not match the code")
	}
	if req.Resource != "" && req.Resource != code.Resource {
		return model.TokenRes{}, model.NewOAuthError("invalid_target", "resource does not match the code")
	}

	// PKCE: SHA256(code_verifier) должен совпасть с code_challenge
	sum := sha256.Sum256([]byte(req.CodeVerifier))
	challenge := base64.RawURLEncoding.EncodeToString(sum[:])
	if req.CodeVerifier == "" || subtle.ConstantTimeCompare([]byte(challenge), []byte(code.CodeChallenge)) != 1 {
		return model.TokenRes{}, model.NewOAuthError("invalid_grant", "code_verifier does not match code_challenge")
	}

	// Имя клиента — в название устройства, чтобы его было видно в списке сессий
	client, _, err := s.oauthRepository.GetClient(ctx, code.ClientID)
	if err != nil {
		return model.TokenRes{}, err
	}

	// Устройство на каждую выдачу: отзыв — удаление устройства, как у телефона
	raw := make([]byte, 8)
	if _, err = rand.Read(raw); err != nil {
		return model.TokenRes{}, errors.InternalServer.Wrap(err)
	}
	tokens, err := s.authService.IssueTokens(ctx, authModel.IssueTokensReq{
		UserID:      code.UserID,
		DeviceID:    "oauth-" + hex.EncodeToString(raw),
		Device:      s.deviceInformation(client.Name),
		Application: s.applicationInformation(),
	})
	if err != nil {
		return model.TokenRes{}, err
	}

	return s.tokenRes(tokens), nil
}

// deviceInformation — устройство OAuth-клиента: без локальной копии данных
func (s *OAuthService) deviceInformation(clientName string) userModel.DeviceInformation {
	return userModel.DeviceInformation{ //nolint:exhaustruct
		NameOS:     osType.Server,
		VersionOS:  s.config.Version,
		DeviceName: clientName,
		ModelName:  "mcp-oauth",
		Stateless:  true,
	}
}

// applicationInformation — приложение OAuth-клиента для записи устройства
func (s *OAuthService) applicationInformation() userModel.ApplicationInformation {
	return userModel.ApplicationInformation{BundleID: "tasks.mcp.oauth", Version: s.config.Version, Build: s.config.Version}
}

// tokenRes собирает ответ эндпоинта токенов
func (s *OAuthService) tokenRes(tokens authModel.Tokens) model.TokenRes {
	return model.TokenRes{
		AccessToken:  tokens.AccessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.config.AccessTokenTTL.Seconds()),
		RefreshToken: tokens.RefreshToken,
	}
}
