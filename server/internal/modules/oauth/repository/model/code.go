package model

import (
	"time"

	"github.com/google/uuid"

	oauthModel "server/internal/modules/oauth/model"
)

// Code — репозиторная модель authorization code
type Code struct {
	Code          string     `db:"code"`           // Случайный код
	ClientID      uuid.UUID  `db:"client_id"`      // Клиент
	UserID        uuid.UUID  `db:"user_id"`        // Пользователь
	RedirectURI   string     `db:"redirect_uri"`   // redirect_uri запроса
	CodeChallenge string     `db:"code_challenge"` // PKCE S256
	Resource      string     `db:"resource"`       // Канонический URL MCP
	ExpiresAt     time.Time  `db:"expires_at"`     // Срок жизни
	UsedAt        *time.Time `db:"used_at"`        // Момент обмена
}

// ConvertToModel конвертирует репозиторную модель в сервисную
func (c Code) ConvertToModel() oauthModel.Code {
	return oauthModel.Code{
		Code:          c.Code,
		ClientID:      c.ClientID,
		UserID:        c.UserID,
		RedirectURI:   c.RedirectURI,
		CodeChallenge: c.CodeChallenge,
		Resource:      c.Resource,
		ExpiresAt:     c.ExpiresAt,
		UsedAt:        c.UsedAt,
	}
}
