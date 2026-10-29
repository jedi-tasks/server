package priorityDDL

import "server/internal/ddl"

const (
	Table          = ddl.SchemaTasks + "." + "priorities"
	TableWithAlias = Table + " " + alias
	alias          = "p"
)

const (
	ColumnID        = "id"
	ColumnUserID    = "user_id"
	ColumnName      = "name"
	ColumnColorHex  = "color_hex"
	ColumnMarks     = "marks"
	ColumnSortOrder = "sort_order"
	ColumnUpdatedAt = "updated_at"
	ColumnDeletedAt = "deleted_at"
)

func WithPrefix(column string) string {
	return alias + "." + column
}
