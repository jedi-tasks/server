package oauthClientDDL

import "server/internal/ddl"

const (
	Table = ddl.SchemaTasks + "." + "oauth_clients"
)

const (
	ColumnID           = "id"            // client_id
	ColumnName         = "name"          // Название клиента
	ColumnRedirectURIs = "redirect_uris" // Разрешенные redirect_uri, JSON-массив
	ColumnCreatedAt    = "created_at"    // Дата регистрации
	ColumnLastUsedAt   = "last_used_at"  // Последняя авторизация
)
