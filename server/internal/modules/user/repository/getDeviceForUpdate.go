package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/google/uuid"

	"pkg/ddlHelper"
	userModel "server/internal/modules/user/model"
	"server/internal/modules/user/repository/deviceDDL"
)

// GetDeviceForUpdate возвращает устройство пользователя с блокировкой строки, должен вызываться внутри транзакции
func (r *UserRepository) GetDeviceForUpdate(ctx context.Context, userID uuid.UUID, deviceID string) (device userModel.Device, err error) {
	ctx, span := tracer.Start(ctx, "GetDeviceForUpdate")
	defer span.End()

	// Блокируем строку устройства для последующего атомарного изменения чекпоинта
	err = r.db.Get(ctx, &device, sq.
		Select(ddlHelper.SelectAll).
		From(deviceDDL.Table).
		Where(sq.Eq{
			deviceDDL.ColumnUserID:   userID,
			deviceDDL.ColumnDeviceID: deviceID,
		}).
		Suffix("FOR UPDATE"),
	)

	return device, err
}
