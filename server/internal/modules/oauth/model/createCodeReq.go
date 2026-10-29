package model

import (
	"time"

	"github.com/google/uuid"
)

// CreateCodeReq — выдача authorization code после входа пользователя
type CreateCodeReq struct {
	Code          string    `validate:"required"` // Случайный код
	ClientID      uuid.UUID `validate:"required"` // Клиент
	UserID        uuid.UUID `validate:"required"` // Пользователь
	RedirectURI   string    `validate:"required"` // redirect_uri запроса
	CodeChallenge string    `validate:"required"` // PKCE S256
	Resource      string    `validate:"required"` // Канонический URL MCP
	ExpiresAt     time.Time `validate:"required"` // Срок жизни
}
