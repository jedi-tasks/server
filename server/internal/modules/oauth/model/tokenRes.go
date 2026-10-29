package model

// TokenRes — ответ эндпоинта токенов (RFC 6749 §5.1)
type TokenRes struct {
	AccessToken  string `json:"access_token"`  // JWT доступа
	TokenType    string `json:"token_type"`    // Bearer
	ExpiresIn    int64  `json:"expires_in"`    // Срок жизни access-токена в секундах
	RefreshToken string `json:"refresh_token"` // JWT обновления
}
