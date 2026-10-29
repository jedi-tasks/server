package service

import (
	"context"

	"github.com/google/uuid"
)

// SetDevicePendingSync записывает чекпоинт и токен последнего вызова Sync, ожидающий подтверждения от устройства
func (s *UserService) SetDevicePendingSync(ctx context.Context, userID uuid.UUID, deviceID string, pendingCheckpoint uint32, pendingSyncToken uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "SetDevicePendingSync")
	defer span.End()

	return s.userRepository.SetDevicePendingSync(ctx, userID, deviceID, pendingCheckpoint, pendingSyncToken)
}
