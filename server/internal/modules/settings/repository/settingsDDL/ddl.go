package settingsDDL

import "server/internal/ddl"

const (
	Table = ddl.SchemaTasks + "." + "settings"
)

const (
	ColumnUserID               = "user_id"
	ColumnDailyCapacityMinutes = "daily_capacity_minutes"
	ColumnTodaySortMode        = "today_sort_mode"
	ColumnUpdatedAt            = "updated_at"
)
