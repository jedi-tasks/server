package grpc

import (
	"context"

	"github.com/google/uuid"

	"server/internal/modules/settings/model"
	settingsService "server/internal/modules/settings/service"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

var _ SettingsService = new(settingsService.SettingsService)

type SettingsService interface {
	GetSettings(ctx context.Context, userID uuid.UUID) (model.Settings, error)
	UpdateSettings(context.Context, model.UpdateSettingsReq) error
	GetVersion(context.Context, model.GetVersionReq) (model.GetVersionRes, error)
}

var _ proto.SettingsEndpointServer = new(SettingsEndpoint)

type SettingsEndpoint struct {
	proto.UnsafeSettingsEndpointServer
	settingsService SettingsService
}

func NewSettingsEndpoint(settingsService SettingsService) *SettingsEndpoint {
	return &SettingsEndpoint{
		settingsService: settingsService,
	}
}
