package model

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// Tag — тег
type Tag struct {
	ID        uuid.UUID `json:"id" db:"id"`                // Идентификатор
	UserID    uuid.UUID `json:"userID" db:"user_id"`       // Идентификатор пользователя-владельца
	Name      string    `json:"name" db:"name"`            // Название тега
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"` // Дата и время последнего изменения
}

// ConvertToProto преобразует Tag в proto-формат
func (m Tag) ConvertToProto() (*proto.Tag, error) {
	return &proto.Tag{
		Id:        m.ID[:],
		Name:      m.Name,
		UpdatedAt: timestamppb.New(m.UpdatedAt),
	}, nil
}
