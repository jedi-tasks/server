package model

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// TaskToTag — связь задачи и тега
type TaskToTag struct {
	UserID    uuid.UUID `json:"userID" db:"user_id"`       // Идентификатор пользователя-владельца
	TaskID    uuid.UUID `json:"taskID" db:"task_id"`       // Идентификатор задачи
	TagID     uuid.UUID `json:"tagID" db:"tag_id"`         // Идентификатор тега
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"` // Дата и время последнего изменения
}

// EntityID возвращает идентификатор связи для аудит-лога
func (m TaskToTag) EntityID() string {
	return TaskToTagEntityID(m.TaskID, m.TagID)
}

// TaskToTagEntityID собирает идентификатор связи для аудит-лога: "taskID:tagID"
func TaskToTagEntityID(taskID, tagID uuid.UUID) string {
	return taskID.String() + ":" + tagID.String()
}

// ConvertToProto преобразует TaskToTag в proto-формат
func (m TaskToTag) ConvertToProto() *proto.TaskToTag {
	return &proto.TaskToTag{
		TaskID:    m.TaskID[:],
		TagID:     m.TagID[:],
		UpdatedAt: timestamppb.New(m.UpdatedAt),
	}
}
