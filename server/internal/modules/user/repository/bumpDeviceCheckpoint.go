package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/google/uuid"

	"server/internal/modules/user/repository/deviceDDL"
)

// BumpDeviceCheckpoint атомарно поднимает чекпоинт устройства до auditLogID, если он больше текущего.
// Используется сразу после собственной мутации устройства, чтобы не требовать от него лишнего Sync
func (r *UserRepository) BumpDeviceCheckpoint(ctx context.Context, userID uuid.UUID, deviceID string, auditLogID uint32) error {
	ctx, span := tracer.Start(ctx, "BumpDeviceCheckpoint")
	defer span.End()

	return r.db.Exec(ctx, sq.Update(deviceDDL.Table).
		Set(deviceDDL.ColumnLastSyncedAuditLogID, sq.Expr("GREATEST("+deviceDDL.ColumnLastSyncedAuditLogID+", ?)", auditLogID)).
		Where(sq.Eq{
			deviceDDL.ColumnUserID:   userID,
			deviceDDL.ColumnDeviceID: deviceID,
		}),
	)
}
