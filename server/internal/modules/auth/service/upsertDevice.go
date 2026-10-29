package service

import (
	"context"

	"github.com/google/uuid"

	userModel "server/internal/modules/user/model"
	userRepoModel "server/internal/modules/user/repository/model"
)

func (s *AuthService) upsertDevice(ctx context.Context, device userModel.Device) error {
	ctx, span := tracer.Start(ctx, "upsertDevice")
	defer span.End()

	// Получаем девайс пользователя
	devices, err := s.userRepository.GetDevices(ctx, userRepoModel.GetDevicesReq{ //nolint:exhaustruct
		DeviceIDs: []string{device.DeviceID},
		UserIDs:   []uuid.UUID{device.UserID},
	})
	if err != nil {
		return err
	}

	// Если девайс не нашелся
	if len(devices) == 0 {

		device.ID = uuid.New()
		err = s.userRepository.CreateDevice(ctx, device)

	} else { // Если девайс нашелся

		// Обновляем у него токен
		err = s.userRepository.UpdateDevice(ctx, userRepoModel.UpdateDeviceReq{
			DeviceID:          device.DeviceID,
			UserID:            device.UserID,
			RefreshToken:      &device.RefreshToken,
			NotificationToken: nil,
			ApplicationInformation: userRepoModel.UpdateApplicationInformationReq{
				BundleID: &device.ApplicationInformation.BundleID,
				Version:  &device.ApplicationInformation.Version,
				Build:    &device.ApplicationInformation.Build,
			},
			DeviceInformation: userRepoModel.UpdateDeviceInformationReq{
				VersionOS: &device.DeviceInformation.VersionOS,
				UserAgent: &device.DeviceInformation.UserAgent,
				IPAddress: &device.DeviceInformation.IPAddress,
			},
		})
	}

	if err != nil {
		return err
	}

	return nil
}
