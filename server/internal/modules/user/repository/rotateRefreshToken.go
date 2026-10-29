package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/google/uuid"

	"server/internal/modules/user/repository/deviceDDL"
)

// RotateRefreshToken атомарно заменяет refresh-токен устройства на новый, сохраняя предыдущий
// токен как валидный ещё на graceWindow - на случай, если ответ с новой парой не дойдёт до клиента
// по сети и он повторит запрос со старым (уже неактуальным на сервере) токеном
func (r *UserRepository) RotateRefreshToken(ctx context.Context, userID uuid.UUID, deviceID string, newRefreshToken string, graceWindow time.Duration) error {
	ctx, span := tracer.Start(ctx, "RotateRefreshToken")
	defer span.End()

	return r.db.Exec(ctx, sq.Update(deviceDDL.Table).
		Set(deviceDDL.ColumnPreviousRefreshToken, sq.Expr(deviceDDL.ColumnRefreshToken)).
		Set(deviceDDL.ColumnPreviousRefreshTokenExpiresAt, time.Now().Add(graceWindow)).
		Set(deviceDDL.ColumnRefreshToken, newRefreshToken).
		Where(sq.Eq{
			deviceDDL.ColumnUserID:   userID,
			deviceDDL.ColumnDeviceID: deviceID,
		}),
	)
}
