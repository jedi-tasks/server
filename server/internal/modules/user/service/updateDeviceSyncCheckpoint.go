package service

import (
	"context"

	"github.com/google/uuid"

	userRepoModel "server/internal/modules/user/repository/model"
)

// UpdateDeviceSyncCheckpoint обновляет идентификатор последней синхронизированной записи аудит-лога для устройства
func (s *UserService) UpdateDeviceSyncCheckpoint(ctx context.Context, userID uuid.UUID, deviceID string, lastSyncedAuditLogID uint32) error {
	ctx, span := tracer.Start(ctx, "UpdateDeviceSyncCheckpoint")
	defer span.End()

	return s.userRepository.UpdateDevice(ctx, userRepoModel.UpdateDeviceReq{ //nolint:exhaustruct
		UserID:               userID,
		DeviceID:             deviceID,
		LastSyncedAuditLogID: &lastSyncedAuditLogID,
	})
}
