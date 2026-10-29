package model

import "github.com/google/uuid"

// AuthorizeReq — параметры запроса авторизации (RFC 6749 §4.1.1 + PKCE + resource)
type AuthorizeReq struct {
	ClientID            uuid.UUID // client_id
	RedirectURI         string    // redirect_uri, должен быть в списке клиента
	ResponseType        string    // Только code
	CodeChallenge       string    // PKCE challenge
	CodeChallengeMethod string    // Только S256
	Resource            string    // Канонический URL MCP, пусто — подставляется свой
	State               string    // Возвращается клиенту как есть
}
