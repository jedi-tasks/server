package model

import (
	"time"

	"github.com/google/uuid"
)

// Client — OAuth-клиент, зарегистрированный динамически. Секрета нет: защита — вход пользователя
type Client struct {
	ID           uuid.UUID // client_id
	Name         string    // Название клиента, показывается на форме входа
	RedirectURIs []string  // Разрешенные redirect_uri
	CreatedAt    time.Time // Дата регистрации
	LastUsedAt   time.Time // Последняя авторизация
}
