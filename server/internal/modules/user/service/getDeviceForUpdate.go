package service

import (
	"context"

	"github.com/google/uuid"

	userModel "server/internal/modules/user/model"
)

// GetDeviceSyncStateForUpdate блокирует строку устройства и возвращает его состояние синхронизации:
// чекпоинт и признак устройства без локальной копии, должен вызываться внутри транзакции
func (s *UserService) GetDeviceSyncStateForUpdate(ctx context.Context, userID uuid.UUID, deviceID string) (userModel.DeviceSyncState, error) {
	ctx, span := tracer.Start(ctx, "GetDeviceSyncStateForUpdate")
	defer span.End()

	device, err := s.userRepository.GetDeviceForUpdate(ctx, userID, deviceID)
	if err != nil {
		return userModel.DeviceSyncState{}, err
	}

	return userModel.DeviceSyncState{
		LastSyncedAuditLogID: device.LastSyncedAuditLogID,
		IsStateless:          device.IsStateless,
	}, nil
}
