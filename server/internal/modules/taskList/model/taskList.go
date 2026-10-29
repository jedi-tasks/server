package model

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"server/internal/enum/sortMode"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// TaskList — список задач
type TaskList struct {
	ID        uuid.UUID         `json:"id" db:"id"`                // Идентификатор
	UserID    uuid.UUID         `json:"userID" db:"user_id"`       // Идентификатор пользователя-владельца
	Name      string            `json:"name" db:"name"`            // Название списка
	IsInbox   bool              `json:"isInbox" db:"is_inbox"`     // Системный список Inbox
	ColorHex  string            `json:"colorHex" db:"color_hex"`   // Цвет списка
	Icon      string            `json:"icon" db:"icon"`            // SF Symbol иконки списка
	SortOrder float64           `json:"sortOrder" db:"sort_order"` // Порядок списка
	SortMode  sortMode.SortMode `json:"sortMode" db:"sort_mode"`   // Режим сортировки задач в списке
	CreatedAt time.Time         `json:"createdAt" db:"created_at"` // Дата и время создания
	UpdatedAt time.Time         `json:"updatedAt" db:"updated_at"` // Дата и время последнего изменения
}

// ConvertToProto преобразует TaskList в proto-формат
func (m TaskList) ConvertToProto() (*proto.TaskList, error) {
	protoSortMode, err := m.SortMode.ConvertToProto()
	if err != nil {
		return nil, err
	}

	return &proto.TaskList{
		Id:        m.ID[:],
		Name:      m.Name,
		IsInbox:   m.IsInbox,
		ColorHex:  m.ColorHex,
		Icon:      m.Icon,
		SortOrder: m.SortOrder,
		SortMode:  protoSortMode,
		CreatedAt: timestamppb.New(m.CreatedAt),
		UpdatedAt: timestamppb.New(m.UpdatedAt),
	}, nil
}
