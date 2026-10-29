package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/user/repository/deviceDDL"
)

// DeleteDevice Удаляет девайс пользователя
func (r *UserRepository) DeleteDevice(ctx context.Context, userID uuid.UUID, deviceID string) error {
	ctx, span := tracer.Start(ctx, "DeleteDevice")
	defer span.End()

	return r.db.Exec(ctx, sq.
		Delete(deviceDDL.Table).
		Where(sq.Eq{
			deviceDDL.ColumnUserID:   userID,
			deviceDDL.ColumnDeviceID: deviceID,
		}),
	)
}
