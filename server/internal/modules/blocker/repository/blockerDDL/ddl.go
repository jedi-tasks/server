package blockerDDL

import "server/internal/ddl"

const (
	Table          = ddl.SchemaTasks + "." + "blockers"
	TableWithAlias = Table + " " + alias
	alias          = "b"
)

const (
	ColumnID             = "id"
	ColumnUserID         = "user_id"
	ColumnTaskID         = "task_id"
	ColumnBlockingTaskID = "blocking_task_id"
	ColumnReason         = "reason"
	ColumnSinceDay       = "since_day"
	ColumnResolvedAt     = "resolved_at"
	ColumnUpdatedAt      = "updated_at"
	ColumnDeletedAt      = "deleted_at"
)

func WithPrefix(column string) string {
	return alias + "." + column
}
