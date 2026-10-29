package service

import (
	"context"

	"server/internal/modules/auth/model"
)

// SignIn авторизует пользователя и возвращает токены доступа
func (s *AuthService) SignIn(ctx context.Context, loginData model.SignInReq) (accessData model.AuthRes, err error) {
	ctx, span := tracer.Start(ctx, "SignIn")
	defer span.End()

	// Проверяем пароль
	accessData.ID, err = s.VerifyPassword(ctx, loginData.Email, loginData.Password)
	if err != nil {
		return accessData, err
	}

	// Выдаем токены и регистрируем устройство
	accessData.Tokens, err = s.IssueTokens(ctx, model.IssueTokensReq{
		UserID:      accessData.ID,
		DeviceID:    loginData.DeviceID,
		Device:      loginData.Device,
		Application: loginData.Application,
	})
	if err != nil {
		return accessData, err
	}

	return accessData, nil
}
