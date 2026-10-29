package service

import (
	"context"

	"github.com/google/uuid"

	"pkg/slices"

	"server/internal/modules/user/model"
	userRepoModel "server/internal/modules/user/repository/model"
)

// GetDevice возвращает устройство пользователя
func (s *UserService) GetDevice(ctx context.Context, userID uuid.UUID, deviceID string) (model.Device, error) {
	ctx, span := tracer.Start(ctx, "GetDevice")
	defer span.End()

	return slices.FirstWithError(s.userRepository.GetDevices(ctx, userRepoModel.GetDevicesReq{ //nolint:exhaustruct
		UserIDs:   []uuid.UUID{userID},
		DeviceIDs: []string{deviceID},
	}))
}
