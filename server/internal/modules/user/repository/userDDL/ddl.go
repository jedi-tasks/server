package userDDL

import "server/internal/ddl"

const (
	Table          = ddl.SchemaTasks + "." + "users"
	TableWithAlias = Table + " " + alias
	alias          = "u"
)

const (
	ColumnID           = "id"
	ColumnName         = "name"
	ColumnEmail        = "email"
	ColumnPasswordHash = "password_hash"
	ColumnTimeCreate   = "time_create"
	ColumnPasswordSalt = "password_salt"
)

func WithPrefix(column string) string {
	return alias + "." + column
}
