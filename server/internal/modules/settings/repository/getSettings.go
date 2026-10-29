package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"pkg/ddlHelper"
	"server/internal/modules/settings/model"
	"server/internal/modules/settings/repository/settingsDDL"
)

// GetSettings возвращает сохраненные настройки пользователя. found == false, если пользователь их еще не менял
func (r *SettingsRepository) GetSettings(ctx context.Context, userID uuid.UUID) (settings model.Settings, found bool, err error) {
	ctx, span := tracer.Start(ctx, "GetSettings")
	defer span.End()

	var rows []model.Settings
	if err = r.db.Select(ctx, &rows, sq.
		Select(ddlHelper.SelectAll).
		From(settingsDDL.Table).
		Where(sq.Eq{settingsDDL.ColumnUserID: userID}),
	); err != nil || len(rows) == 0 {
		return settings, false, err
	}

	return rows[0], true, nil
}
