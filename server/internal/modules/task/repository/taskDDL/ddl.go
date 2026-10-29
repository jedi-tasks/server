package taskDDL

import "server/internal/ddl"

const (
	Table          = ddl.SchemaTasks + "." + "tasks"
	TableWithAlias = Table + " " + alias
	alias          = "t"
)

const (
	ColumnID              = "id"
	ColumnUserID          = "user_id"
	ColumnListID          = "list_id"
	ColumnSectionID       = "section_id"
	ColumnTitle           = "title"
	ColumnNotes           = "notes"
	ColumnPriorityID      = "priority_id"
	ColumnPlannedDay      = "planned_day"
	ColumnDeadlineDay     = "deadline_day"
	ColumnDurationMinutes = "duration_minutes"
	ColumnManualOrder     = "manual_order"
	ColumnCompletedAt     = "completed_at"
	ColumnExternalID      = "external_id"
	ColumnCreatedAt       = "created_at"
	ColumnUpdatedAt       = "updated_at"
	ColumnDeletedAt       = "deleted_at"
)

func WithPrefix(column string) string {
	return alias + "." + column
}
