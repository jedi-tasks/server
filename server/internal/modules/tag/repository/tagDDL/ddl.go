package tagDDL

import "server/internal/ddl"

const (
	Table          = ddl.SchemaTasks + "." + "tags"
	TableWithAlias = Table + " " + alias
	alias          = "t"
)

const (
	ColumnID        = "id"
	ColumnUserID    = "user_id"
	ColumnName      = "name"
	ColumnUpdatedAt = "updated_at"
	ColumnDeletedAt = "deleted_at"
)

func WithPrefix(column string) string {
	return alias + "." + column
}
