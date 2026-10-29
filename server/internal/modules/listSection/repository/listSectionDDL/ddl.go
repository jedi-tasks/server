package listSectionDDL

import "server/internal/ddl"

const (
	Table          = ddl.SchemaTasks + "." + "list_sections"
	TableWithAlias = Table + " " + alias
	alias          = "ls"
)

const (
	ColumnID        = "id"
	ColumnUserID    = "user_id"
	ColumnListID    = "list_id"
	ColumnName      = "name"
	ColumnSortOrder = "sort_order"
	ColumnUpdatedAt = "updated_at"
	ColumnDeletedAt = "deleted_at"
)

func WithPrefix(column string) string {
	return alias + "." + column
}
