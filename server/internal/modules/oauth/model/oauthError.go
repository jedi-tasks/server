package model

// OAuthError — ошибка протокола OAuth 2.1 (RFC 6749 §5.2): код для клиента и описание
type OAuthError struct {
	Code        string // invalid_request, invalid_client, invalid_grant, unauthorized_client, unsupported_grant_type, access_denied
	Description string // Текст для разработчика клиента
}

// Error возвращает текст ошибки
func (e *OAuthError) Error() string {
	return e.Code + ": " + e.Description
}

// NewOAuthError создает ошибку протокола
func NewOAuthError(code, description string) *OAuthError {
	return &OAuthError{Code: code, Description: description}
}

// ErrInvalidCredentials — неверная почта или пароль на форме входа: клиенту не уходит, форма показывается снова
var ErrInvalidCredentials = NewOAuthError("access_denied", "неверная почта или пароль")
