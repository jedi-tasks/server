package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateBlockerReq struct {
	ID             uuid.UUID  // Идентификатор
	UserID         uuid.UUID  // Идентификатор пользователя-владельца
	TaskID         uuid.UUID  // Заблокированная задача
	BlockingTaskID *uuid.UUID // Блокирующая задача
	Reason         *string    // Текстовая причина
	SinceDay       string     // Дата начала, yyyy-MM-dd
	UpdatedAt      time.Time  // Время создания на клиенте
}
