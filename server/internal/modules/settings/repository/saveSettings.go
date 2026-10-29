package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"server/internal/modules/settings/model"
	"server/internal/modules/settings/repository/settingsDDL"
)

// SaveSettings сохраняет настройки, если изменение новее сохраненного. Возвращает false, если на
// сервере уже более новая версия
func (r *SettingsRepository) SaveSettings(ctx context.Context, settings model.Settings) (saved bool, err error) {
	ctx, span := tracer.Start(ctx, "SaveSettings")
	defer span.End()

	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Insert(settingsDDL.Table).
		SetMap(map[string]any{
			settingsDDL.ColumnUserID:               settings.UserID,
			settingsDDL.ColumnDailyCapacityMinutes: settings.DailyCapacityMinutes,
			settingsDDL.ColumnTodaySortMode:        settings.TodaySortMode,
			settingsDDL.ColumnUpdatedAt:            settings.UpdatedAt,
		}).
		Suffix(
			"ON CONFLICT ("+settingsDDL.ColumnUserID+") DO UPDATE SET "+
				settingsDDL.ColumnDailyCapacityMinutes+" = EXCLUDED."+settingsDDL.ColumnDailyCapacityMinutes+", "+
				settingsDDL.ColumnTodaySortMode+" = EXCLUDED."+settingsDDL.ColumnTodaySortMode+", "+
				settingsDDL.ColumnUpdatedAt+" = EXCLUDED."+settingsDDL.ColumnUpdatedAt+" "+
				"WHERE "+settingsDDL.Table+"."+settingsDDL.ColumnUpdatedAt+" < EXCLUDED."+settingsDDL.ColumnUpdatedAt,
		),
	)
	if err != nil {
		return false, err
	}

	return rows != 0, nil
}
