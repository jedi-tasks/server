package service

import (
	"context"

	"server/internal/modules/auth/model"
	"server/internal/modules/auth/service/utils"
	userModel "server/internal/modules/user/model"
)

// IssueTokens выдает пару токенов проверенному пользователю и создает или обновляет его устройство
func (s *AuthService) IssueTokens(ctx context.Context, req model.IssueTokensReq) (model.Tokens, error) {
	ctx, span := tracer.Start(ctx, "IssueTokens")
	defer span.End()

	// Создаем пару токенов
	tokens, err := utils.CreatePairTokens(ctx, req.UserID, req.DeviceID)
	if err != nil {
		return tokens, err
	}

	// Создаем или обновляем девайс пользователя
	err = s.upsertDevice(ctx, userModel.Device{ //nolint:exhaustruct
		DeviceInformation:      req.Device,
		NotificationToken:      nil,
		RefreshToken:           tokens.RefreshToken,
		UserID:                 req.UserID,
		DeviceID:               req.DeviceID,
		ApplicationInformation: req.Application,
	})
	if err != nil {
		return tokens, err
	}

	return tokens, nil
}
