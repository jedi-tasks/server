package model

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// ListSection — раздел внутри списка
type ListSection struct {
	ID        uuid.UUID `json:"id" db:"id"`                // Идентификатор
	UserID    uuid.UUID `json:"userID" db:"user_id"`       // Идентификатор пользователя-владельца
	ListID    uuid.UUID `json:"listID" db:"list_id"`       // Идентификатор списка
	Name      string    `json:"name" db:"name"`            // Название раздела
	SortOrder float64   `json:"sortOrder" db:"sort_order"` // Порядок раздела
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"` // Дата и время последнего изменения
}

// ConvertToProto преобразует ListSection в proto-формат
func (m ListSection) ConvertToProto() (*proto.ListSection, error) {
	return &proto.ListSection{
		Id:        m.ID[:],
		ListID:    m.ListID[:],
		Name:      m.Name,
		SortOrder: m.SortOrder,
		UpdatedAt: timestamppb.New(m.UpdatedAt),
	}, nil
}
