package service

import (
	"context"

	"github.com/google/uuid"
)

// BumpDeviceCheckpoint атомарно поднимает чекпоинт устройства до auditLogID, если он больше текущего.
// Используется сразу после собственной мутации устройства, чтобы не требовать от него лишнего Sync
func (s *UserService) BumpDeviceCheckpoint(ctx context.Context, userID uuid.UUID, deviceID string, auditLogID uint32) error {
	ctx, span := tracer.Start(ctx, "BumpDeviceCheckpoint")
	defer span.End()

	return s.userRepository.BumpDeviceCheckpoint(ctx, userID, deviceID, auditLogID)
}
