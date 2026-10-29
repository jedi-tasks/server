package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateTaskReq struct {
	ID              uuid.UUID  // Идентификатор
	UserID          uuid.UUID  // Идентификатор пользователя-владельца
	ListID          uuid.UUID  // Идентификатор списка
	SectionID       *uuid.UUID // Идентификатор раздела
	Title           string     // Название задачи
	Notes           string     // Описание
	PriorityID      *uuid.UUID // Идентификатор приоритета
	PlannedDay      *string    // Планируемый срок, yyyy-MM-dd
	DeadlineDay     *string    // Дедлайн, yyyy-MM-dd
	DurationMinutes *int32     // Планируемое время в минутах
	ManualOrder     float64    // Ручной порядок
	ExternalID      *string    // Идентификатор напоминания при импорте
	CreatedAt       time.Time  // Время создания на клиенте
}
