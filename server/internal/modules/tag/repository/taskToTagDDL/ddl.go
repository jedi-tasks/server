package taskToTagDDL

import "server/internal/ddl"

const (
	Table          = ddl.SchemaTasks + "." + "tasks_to_tags"
	TableWithAlias = Table + " " + alias
	alias          = "ttt"
)

const (
	ColumnUserID    = "user_id"
	ColumnTaskID    = "task_id"
	ColumnTagID     = "tag_id"
	ColumnUpdatedAt = "updated_at"
	ColumnDeletedAt = "deleted_at"
)

func WithPrefix(column string) string {
	return alias + "." + column
}
