package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"server/internal/modules/oauth/model"
	"server/internal/utils/errors"
)

// codeTTL — срок жизни authorization code
const codeTTL = 10 * time.Minute

// GrantCode проверяет пароль пользователя и выдает одноразовый код для клиента
func (s *OAuthService) GrantCode(ctx context.Context, req model.GrantCodeReq) (code string, err error) {
	ctx, span := tracer.Start(ctx, "GrantCode")
	defer span.End()

	// Параметры авторизации проверяем повторно: форма могла быть подделана
	if _, err = s.ValidateAuthorize(ctx, req.Authorize); err != nil {
		return "", err
	}

	// Проверяем пароль
	userID, err := s.authService.VerifyPassword(ctx, req.Email, req.Password)
	if err != nil {
		return "", model.ErrInvalidCredentials
	}

	// Выдаем код
	raw := make([]byte, 32)
	if _, err = rand.Read(raw); err != nil {
		return "", errors.InternalServer.Wrap(err)
	}
	code = base64.RawURLEncoding.EncodeToString(raw)
	resource := req.Authorize.Resource
	if resource == "" {
		resource = s.config.PublicURL
	}
	now := time.Now()
	err = s.oauthRepository.CreateCode(ctx, model.CreateCodeReq{
		Code:          code,
		ClientID:      req.Authorize.ClientID,
		UserID:        userID,
		RedirectURI:   req.Authorize.RedirectURI,
		CodeChallenge: req.Authorize.CodeChallenge,
		Resource:      resource,
		ExpiresAt:     now.Add(codeTTL),
	})
	if err != nil {
		return "", err
	}

	// Отмечаем активность клиента
	if err = s.oauthRepository.TouchClient(ctx, req.Authorize.ClientID, now); err != nil {
		return "", err
	}

	return code, nil
}
