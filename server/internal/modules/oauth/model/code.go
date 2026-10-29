package model

import (
	"time"

	"github.com/google/uuid"
)

// Code — одноразовый authorization code с привязкой к клиенту, пользователю и PKCE
type Code struct {
	Code          string     // Случайный код
	ClientID      uuid.UUID  // Клиент, запросивший код
	UserID        uuid.UUID  // Пользователь, вошедший на форме
	RedirectURI   string     // redirect_uri запроса авторизации, должен совпасть при обмене
	CodeChallenge string     // PKCE S256: BASE64URL(SHA256(code_verifier))
	Resource      string     // Канонический URL MCP, для которого выдан код
	ExpiresAt     time.Time  // Срок жизни кода
	UsedAt        *time.Time // Момент обмена, nil — не использован
}
