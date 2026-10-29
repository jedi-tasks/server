package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateListSectionReq struct {
	ID        uuid.UUID // Идентификатор
	UserID    uuid.UUID // Идентификатор пользователя-владельца
	ListID    uuid.UUID // Идентификатор списка
	Name      string    // Название раздела
	SortOrder float64   // Порядок раздела
	UpdatedAt time.Time // Время создания на клиенте
}
