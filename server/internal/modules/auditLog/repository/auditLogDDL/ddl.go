package auditLogDDL

import "server/internal/ddl"

const (
	Table = ddl.SchemaTasks + "." + "audit_log"
)

const (
	ColumnID             = "id"
	ColumnEntity         = "entity"
	ColumnMethod         = "method"
	ColumnEntityID       = "entity_id"
	ColumnSnapshotBefore = "snapshot_before"
	ColumnSnapshotAfter  = "snapshot_after"
	ColumnUserID         = "user_id"
	ColumnDeviceID       = "device_id"
	ColumnDatetimeCreate = "datetime_create"
)
