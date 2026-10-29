package model

// ExchangeCodeReq — обмен кода на токены (grant_type=authorization_code)
type ExchangeCodeReq struct {
	Code         string // Код из редиректа
	ClientID     string // client_id
	RedirectURI  string // Тот же redirect_uri, что в запросе авторизации
	CodeVerifier string // PKCE verifier
	Resource     string // Канонический URL MCP, пусто — не проверяется повторно
}
