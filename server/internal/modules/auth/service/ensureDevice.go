package service

import (
	"context"

	"github.com/google/uuid"

	"pkg/validator"
	"server/internal/modules/auth/model"
	userModel "server/internal/modules/user/model"
	userRepoModel "server/internal/modules/user/repository/model"
)

// EnsureDevice создает устройство пользователя, если его еще нет. Нужен устройствам без входа по
// паролю, например MCP: личность уже подтверждена токеном, а строка устройства нужна для чекпоинта
// и аудит-лога. Возвращает true, если устройство создано этим вызовом
func (s *AuthService) EnsureDevice(ctx context.Context, req model.EnsureDeviceReq) (created bool, err error) {
	ctx, span := tracer.Start(ctx, "EnsureDevice")
	defer span.End()

	// Валидируем запрос
	if err = validator.Validate(req); err != nil {
		return false, err
	}

	// Проверяем, есть ли устройство
	devices, err := s.userRepository.GetDevices(ctx, userRepoModel.GetDevicesReq{ //nolint:exhaustruct
		DeviceIDs: []string{req.DeviceID},
		UserIDs:   []uuid.UUID{req.UserID},
	})
	if err != nil {
		return false, err
	}
	if len(devices) != 0 {
		return false, nil
	}

	// Создаем устройство без refresh-токена: входа по паролю у него не будет
	err = s.userRepository.CreateDevice(ctx, userModel.Device{ //nolint:exhaustruct
		ID:                     uuid.New(),
		DeviceInformation:      req.Device,
		ApplicationInformation: req.Application,
		UserID:                 req.UserID,
		DeviceID:               req.DeviceID,
	})
	if err != nil {
		return false, err
	}

	return true, nil
}
