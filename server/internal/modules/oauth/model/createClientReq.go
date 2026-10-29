package model

import (
	"time"

	"github.com/google/uuid"
)

// CreateClientReq — регистрация клиента
type CreateClientReq struct {
	ID           uuid.UUID `validate:"required"` // client_id, генерируется сервером
	Name         string    // Название клиента
	RedirectURIs []string  `validate:"required"` // Разрешенные redirect_uri
	CreatedAt    time.Time `validate:"required"` // Дата регистрации
}
