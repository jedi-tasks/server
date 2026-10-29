package service

import (
	"context"

	"pkg/validator"

	"server/internal/modules/auth/model"
	userModel "server/internal/modules/user/model"
)

// ensureUserDeviceID — устройство в аудит-логе для пользователя, созданного сервером при запуске
const ensureUserDeviceID = "server"

// EnsureUser создает пользователя, если пользователя с таким email еще нет.
// Существующего пользователя не меняет, в том числе его пароль. Возвращает признак создания
func (s *AuthService) EnsureUser(ctx context.Context, req model.EnsureUserReq) (created bool, err error) {
	ctx, span := tracer.Start(ctx, "EnsureUser")
	defer span.End()

	// Проверяем, что заданы все данные пользователя
	if err = validator.Validate(req); err != nil {
		return false, err
	}

	// Пропускаем, если пользователь уже есть
	users, err := s.userRepository.GetUsers(ctx, userModel.GetUsersReq{Emails: []string{req.Email}}) //nolint:exhaustruct
	if err != nil {
		return false, err
	}
	if len(users) != 0 {
		return false, nil
	}

	// Создаем пользователя с данными по умолчанию
	return true, s.generalRepository.WithinTransaction(ctx, func(ctx context.Context) error {
		_, err := s.createUser(ctx, model.CreateUserReq{
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
			DeviceID: ensureUserDeviceID,
		})
		return err
	})
}
