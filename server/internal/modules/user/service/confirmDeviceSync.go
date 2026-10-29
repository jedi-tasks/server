package service

import (
	"context"

	"github.com/google/uuid"

	"server/internal/utils/errors"
)

// ConfirmDeviceSync принимает от устройства подтверждение, что оно корректно применило ответ Sync,
// и атомарно фиксирует чекпоинт. Если токен устарел или не совпадает с ожидающим - возвращает ошибку
func (s *UserService) ConfirmDeviceSync(ctx context.Context, userID uuid.UUID, deviceID string, pendingSyncToken uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "ConfirmDeviceSync")
	defer span.End()

	rowsAffected, err := s.userRepository.ConfirmDeviceSync(ctx, userID, deviceID, pendingSyncToken)
	if err != nil {
		return err
	}

	// Токен не совпал с текущим ожидающим - устаревшее или повторное подтверждение
	if rowsAffected == 0 {
		return errors.NeedToSync.New("Неактуальный токен синхронизации").
			WithContextParams(ctx)
	}

	return nil
}
