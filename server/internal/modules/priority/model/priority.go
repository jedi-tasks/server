package model

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// Priority — уровень приоритета из справочника
type Priority struct {
	ID        uuid.UUID `json:"id" db:"id"`                // Идентификатор
	UserID    uuid.UUID `json:"userID" db:"user_id"`       // Идентификатор пользователя-владельца
	Name      string    `json:"name" db:"name"`            // Название
	ColorHex  string    `json:"colorHex" db:"color_hex"`   // Цвет
	Marks     string    `json:"marks" db:"marks"`          // Знаки в строке задачи
	SortOrder float64   `json:"sortOrder" db:"sort_order"` // Порядок важности
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"` // Дата и время последнего изменения
}

// ConvertToProto преобразует Priority в proto-формат
func (m Priority) ConvertToProto() (*proto.Priority, error) {
	return &proto.Priority{
		Id:        m.ID[:],
		Name:      m.Name,
		ColorHex:  m.ColorHex,
		Marks:     m.Marks,
		SortOrder: m.SortOrder,
		UpdatedAt: timestamppb.New(m.UpdatedAt),
	}, nil
}
