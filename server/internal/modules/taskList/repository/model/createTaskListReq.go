package model

import (
	"time"

	"github.com/google/uuid"

	"server/internal/enum/sortMode"
)

type CreateTaskListReq struct {
	ID        uuid.UUID         // Идентификатор
	UserID    uuid.UUID         // Идентификатор пользователя-владельца
	Name      string            // Название списка
	ColorHex  string            // Цвет списка
	Icon      string            // SF Symbol иконки списка
	SortOrder float64           // Порядок списка
	SortMode  sortMode.SortMode // Режим сортировки задач в списке
	CreatedAt time.Time         // Время создания на клиенте
}
