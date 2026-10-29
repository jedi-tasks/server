package model

import (
	"time"

	"github.com/google/uuid"
)

type CreatePriorityReq struct {
	ID        uuid.UUID // Идентификатор
	UserID    uuid.UUID // Идентификатор пользователя-владельца
	Name      string    // Название
	ColorHex  string    // Цвет
	Marks     string    // Знаки в строке задачи
	SortOrder float64   // Порядок важности
	UpdatedAt time.Time // Время создания на клиенте
}
