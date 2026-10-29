package oauthCodeDDL

import "server/internal/ddl"

const (
	Table = ddl.SchemaTasks + "." + "oauth_codes"
)

const (
	ColumnCode          = "code"           // Случайный код
	ColumnClientID      = "client_id"      // Клиент, запросивший код
	ColumnUserID        = "user_id"        // Пользователь, вошедший на форме
	ColumnRedirectURI   = "redirect_uri"   // redirect_uri запроса авторизации
	ColumnCodeChallenge = "code_challenge" // PKCE S256
	ColumnResource      = "resource"       // Канонический URL MCP
	ColumnExpiresAt     = "expires_at"     // Срок жизни кода
	ColumnUsedAt        = "used_at"        // Момент обмена, NULL — не использован
)
