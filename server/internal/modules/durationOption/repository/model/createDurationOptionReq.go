package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateDurationOptionReq struct {
	ID        uuid.UUID // Идентификатор
	UserID    uuid.UUID // Идентификатор пользователя-владельца
	Title     string    // Подпись
	Minutes   int32     // Минуты, меньше нуля — непрогнозируемо долго
	UpdatedAt time.Time // Время создания на клиенте
}
