package model

// RegisterClientRes — ответ регистрации (RFC 7591 §3.2.1)
type RegisterClientRes struct {
	ClientID                string   `json:"client_id"`                  // Выданный идентификатор
	ClientName              string   `json:"client_name"`                // Название клиента
	RedirectURIs            []string `json:"redirect_uris"`              // Разрешенные redirect_uri
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method"` // none — публичный клиент без секрета
	GrantTypes              []string `json:"grant_types"`                // authorization_code, refresh_token
	ResponseTypes           []string `json:"response_types"`             // code
}
