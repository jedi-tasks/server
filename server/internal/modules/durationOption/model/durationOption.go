package model

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// DurationOption — значение из справочника длительностей
type DurationOption struct {
	ID        uuid.UUID `json:"id" db:"id"`                // Идентификатор
	UserID    uuid.UUID `json:"userID" db:"user_id"`       // Идентификатор пользователя-владельца
	Title     string    `json:"title" db:"title"`          // Подпись
	Minutes   int32     `json:"minutes" db:"minutes"`      // Минуты, меньше нуля — непрогнозируемо долго
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"` // Дата и время последнего изменения
}

// ConvertToProto преобразует DurationOption в proto-формат
func (m DurationOption) ConvertToProto() (*proto.DurationOption, error) {
	return &proto.DurationOption{
		Id:        m.ID[:],
		Title:     m.Title,
		Minutes:   m.Minutes,
		UpdatedAt: timestamppb.New(m.UpdatedAt),
	}, nil
}
