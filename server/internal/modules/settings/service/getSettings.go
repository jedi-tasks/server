package service

import (
	"context"

	"github.com/google/uuid"

	"server/internal/modules/settings/model"
)

// GetSettings возвращает настройки пользователя, а если он их еще не менял — настройки по умолчанию
func (s *SettingsService) GetSettings(ctx context.Context, userID uuid.UUID) (model.Settings, error) {
	ctx, span := tracer.Start(ctx, "GetSettings")
	defer span.End()

	settings, found, err := s.settingsRepository.GetSettings(ctx, userID)
	if err != nil {
		return model.Settings{}, err
	}
	if !found {
		return model.DefaultSettings(userID), nil
	}

	return settings, nil
}
