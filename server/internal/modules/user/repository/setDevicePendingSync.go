package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/google/uuid"

	"server/internal/modules/user/repository/deviceDDL"
)

// SetDevicePendingSync записывает чекпоинт и токен последнего вызова Sync, ожидающий подтверждения от устройства
func (r *UserRepository) SetDevicePendingSync(ctx context.Context, userID uuid.UUID, deviceID string, pendingCheckpoint uint32, pendingSyncToken uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "SetDevicePendingSync")
	defer span.End()

	return r.db.Exec(ctx, sq.Update(deviceDDL.Table).
		SetMap(map[string]any{
			deviceDDL.ColumnPendingCheckpoint: pendingCheckpoint,
			deviceDDL.ColumnPendingSyncToken:  pendingSyncToken,
		}).
		Where(sq.Eq{
			deviceDDL.ColumnUserID:   userID,
			deviceDDL.ColumnDeviceID: deviceID,
		}),
	)
}
