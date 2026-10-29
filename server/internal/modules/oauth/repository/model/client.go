package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	oauthModel "server/internal/modules/oauth/model"
	"server/internal/utils/errors"
)

// Client — репозиторная модель клиента, redirect_uri хранятся JSON-массивом
type Client struct {
	ID           uuid.UUID       `db:"id"`            // client_id
	Name         string          `db:"name"`          // Название клиента
	RedirectURIs json.RawMessage `db:"redirect_uris"` // JSON-массив строк
	CreatedAt    time.Time       `db:"created_at"`    // Дата регистрации
	LastUsedAt   time.Time       `db:"last_used_at"`  // Последняя авторизация
}

// ConvertToModel конвертирует репозиторную модель в сервисную
func (c Client) ConvertToModel() (oauthModel.Client, error) {
	var redirectURIs []string
	if err := json.Unmarshal(c.RedirectURIs, &redirectURIs); err != nil {
		return oauthModel.Client{}, errors.InternalServer.Wrap(err)
	}
	return oauthModel.Client{
		ID:           c.ID,
		Name:         c.Name,
		RedirectURIs: redirectURIs,
		CreatedAt:    c.CreatedAt,
		LastUsedAt:   c.LastUsedAt,
	}, nil
}
