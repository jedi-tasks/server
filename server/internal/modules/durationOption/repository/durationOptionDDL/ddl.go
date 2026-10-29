package durationOptionDDL

import "server/internal/ddl"

const (
	Table          = ddl.SchemaTasks + "." + "duration_options"
	TableWithAlias = Table + " " + alias
	alias          = "do"
)

const (
	ColumnID        = "id"
	ColumnUserID    = "user_id"
	ColumnTitle     = "title"
	ColumnMinutes   = "minutes"
	ColumnUpdatedAt = "updated_at"
	ColumnDeletedAt = "deleted_at"
)

func WithPrefix(column string) string {
	return alias + "." + column
}
