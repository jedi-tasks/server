package service

import (
	"context"

	settingsModel "server/internal/modules/settings/model"
)

// GetVersion возвращает версию сервера
func (s *SettingsService) GetVersion(ctx context.Context, _ settingsModel.GetVersionReq) (res settingsModel.GetVersionRes, err error) {
	_, span := tracer.Start(ctx, "GetVersion")
	defer span.End()

	return settingsModel.GetVersionRes{
		Version: settingsModel.Version{
			Version: s.version.Version,
			Build:   s.version.Build,
		},
	}, nil
}
