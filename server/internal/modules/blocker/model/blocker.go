package model

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// Blocker — блокер задачи
type Blocker struct {
	ID             uuid.UUID  `json:"id" db:"id"`                           // Идентификатор
	UserID         uuid.UUID  `json:"userID" db:"user_id"`                  // Идентификатор пользователя-владельца
	TaskID         uuid.UUID  `json:"taskID" db:"task_id"`                  // Заблокированная задача
	BlockingTaskID *uuid.UUID `json:"blockingTaskID" db:"blocking_task_id"` // Блокирующая задача
	Reason         *string    `json:"reason" db:"reason"`                   // Текстовая причина
	SinceDay       string     `json:"sinceDay" db:"since_day"`              // Дата начала, yyyy-MM-dd
	ResolvedAt     *time.Time `json:"resolvedAt" db:"resolved_at"`          // Дата и время снятия
	UpdatedAt      time.Time  `json:"updatedAt" db:"updated_at"`            // Дата и время последнего изменения
}

// ConvertToProto преобразует Blocker в proto-формат
func (m Blocker) ConvertToProto() (*proto.Blocker, error) {
	var blockingTaskID []byte
	if m.BlockingTaskID != nil {
		blockingTaskID = m.BlockingTaskID[:]
	}
	var resolvedAt *timestamppb.Timestamp
	if m.ResolvedAt != nil {
		resolvedAt = timestamppb.New(*m.ResolvedAt)
	}

	return &proto.Blocker{
		Id:             m.ID[:],
		TaskID:         m.TaskID[:],
		BlockingTaskID: blockingTaskID,
		Reason:         m.Reason,
		SinceDay:       m.SinceDay,
		ResolvedAt:     resolvedAt,
		UpdatedAt:      timestamppb.New(m.UpdatedAt),
	}, nil
}
