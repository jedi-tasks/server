package service

import (
	"context"

	"github.com/google/uuid"
)

// GetDeviceLastSyncedAuditLogIDForUpdate блокирует строку устройства и возвращает его текущий
// чекпоинт синхронизации, должен вызываться внутри транзакции
func (s *UserService) GetDeviceLastSyncedAuditLogIDForUpdate(ctx context.Context, userID uuid.UUID, deviceID string) (uint32, error) {
	ctx, span := tracer.Start(ctx, "GetDeviceLastSyncedAuditLogIDForUpdate")
	defer span.End()

	device, err := s.userRepository.GetDeviceForUpdate(ctx, userID, deviceID)
	if err != nil {
		return 0, err
	}

	return device.LastSyncedAuditLogID, nil
}
