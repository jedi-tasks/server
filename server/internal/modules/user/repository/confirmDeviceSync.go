package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/google/uuid"

	"server/internal/modules/user/repository/deviceDDL"
)

// ConfirmDeviceSync атомарно принимает подтверждение синхронизации по опорному токену: переносит
// pendingCheckpoint в last_synced_audit_log_id и очищает поля ожидающей синхронизации. Возвращает
// число затронутых строк - 0 означает, что токен не совпал с текущим ожидающим (устаревший/повторный вызов)
func (r *UserRepository) ConfirmDeviceSync(ctx context.Context, userID uuid.UUID, deviceID string, pendingSyncToken uuid.UUID) (uint32, error) {
	ctx, span := tracer.Start(ctx, "ConfirmDeviceSync")
	defer span.End()

	return r.db.ExecWithRowsAffected(ctx, sq.Update(deviceDDL.Table).
		Set(deviceDDL.ColumnLastSyncedAuditLogID, sq.Expr("GREATEST("+deviceDDL.ColumnLastSyncedAuditLogID+", "+deviceDDL.ColumnPendingCheckpoint+")")).
		Set(deviceDDL.ColumnPendingCheckpoint, nil).
		Set(deviceDDL.ColumnPendingSyncToken, nil).
		Where(sq.Eq{
			deviceDDL.ColumnUserID:           userID,
			deviceDDL.ColumnDeviceID:         deviceID,
			deviceDDL.ColumnPendingSyncToken: pendingSyncToken,
		}),
	)
}
