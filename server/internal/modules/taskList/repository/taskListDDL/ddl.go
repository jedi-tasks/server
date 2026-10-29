package taskListDDL

import "server/internal/ddl"

const (
	Table          = ddl.SchemaTasks + "." + "task_lists"
	TableWithAlias = Table + " " + alias
	alias          = "tl"
)

const (
	ColumnID        = "id"
	ColumnUserID    = "user_id"
	ColumnName      = "name"
	ColumnIsInbox   = "is_inbox"
	ColumnColorHex  = "color_hex"
	ColumnIcon      = "icon"
	ColumnSortOrder = "sort_order"
	ColumnSortMode  = "sort_mode"
	ColumnCreatedAt = "created_at"
	ColumnUpdatedAt = "updated_at"
	ColumnDeletedAt = "deleted_at"
)

func WithPrefix(column string) string {
	return alias + "." + column
}
