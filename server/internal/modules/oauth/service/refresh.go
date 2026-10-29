package service

import (
	"context"

	authModel "server/internal/modules/auth/model"
	"server/internal/modules/oauth/model"
	"server/internal/utils/errors"
)

// Refresh обменивает refresh-токен на новую пару через обычную ротацию устройства
func (s *OAuthService) Refresh(ctx context.Context, refreshToken string) (model.TokenRes, error) {
	ctx, span := tracer.Start(ctx, "Refresh")
	defer span.End()

	res, err := s.authService.RefreshTokens(ctx, authModel.RefreshTokensReq{
		Token:       refreshToken,
		Application: s.applicationInformation(),
		Device:      s.deviceInformation(""),
	})
	if err != nil {
		// Отозванная или просроченная сессия для клиента — invalid_grant, он начнет авторизацию заново
		if errors.GetCategory(errors.CastError(err).ErrorType) == errors.CategoryNeedToLogout {
			return model.TokenRes{}, model.NewOAuthError("invalid_grant", "refresh token is invalid or revoked")
		}
		return model.TokenRes{}, err
	}

	return s.tokenRes(res.Tokens), nil
}
