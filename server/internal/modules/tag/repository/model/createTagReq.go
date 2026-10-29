package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateTagReq struct {
	ID        uuid.UUID // Идентификатор
	UserID    uuid.UUID // Идентификатор пользователя-владельца
	Name      string    // Название тега
	UpdatedAt time.Time // Время создания на клиенте
}
