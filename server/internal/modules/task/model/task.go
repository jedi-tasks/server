package model

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// Task — задача
type Task struct {
	ID              uuid.UUID  `json:"id" db:"id"`                            // Идентификатор
	UserID          uuid.UUID  `json:"userID" db:"user_id"`                   // Идентификатор пользователя-владельца
	ListID          uuid.UUID  `json:"listID" db:"list_id"`                   // Идентификатор списка
	SectionID       *uuid.UUID `json:"sectionID" db:"section_id"`             // Идентификатор раздела
	Title           string     `json:"title" db:"title"`                      // Название задачи
	Notes           string     `json:"notes" db:"notes"`                      // Описание
	PriorityID      *uuid.UUID `json:"priorityID" db:"priority_id"`           // Идентификатор приоритета
	PlannedDay      *string    `json:"plannedDay" db:"planned_day"`           // Планируемый срок, yyyy-MM-dd
	DeadlineDay     *string    `json:"deadlineDay" db:"deadline_day"`         // Дедлайн, yyyy-MM-dd
	DurationMinutes *int32     `json:"durationMinutes" db:"duration_minutes"` // Планируемое время в минутах
	ManualOrder     float64    `json:"manualOrder" db:"manual_order"`         // Ручной порядок
	CompletedAt     *time.Time `json:"completedAt" db:"completed_at"`         // Дата и время выполнения
	ExternalID      *string    `json:"externalID" db:"external_id"`           // Идентификатор напоминания при импорте
	CreatedAt       time.Time  `json:"createdAt" db:"created_at"`             // Дата и время создания
	UpdatedAt       time.Time  `json:"updatedAt" db:"updated_at"`             // Дата и время последнего изменения
}

// ConvertToProto преобразует Task в proto-формат
func (m Task) ConvertToProto() (*proto.Task, error) {
	var sectionID []byte
	if m.SectionID != nil {
		sectionID = m.SectionID[:]
	}
	var priorityID []byte
	if m.PriorityID != nil {
		priorityID = m.PriorityID[:]
	}
	var completedAt *timestamppb.Timestamp
	if m.CompletedAt != nil {
		completedAt = timestamppb.New(*m.CompletedAt)
	}

	return &proto.Task{
		Id:              m.ID[:],
		ListID:          m.ListID[:],
		SectionID:       sectionID,
		Title:           m.Title,
		Notes:           m.Notes,
		PriorityID:      priorityID,
		PlannedDay:      m.PlannedDay,
		DeadlineDay:     m.DeadlineDay,
		DurationMinutes: m.DurationMinutes,
		ManualOrder:     m.ManualOrder,
		CompletedAt:     completedAt,
		ExternalID:      m.ExternalID,
		CreatedAt:       timestamppb.New(m.CreatedAt),
		UpdatedAt:       timestamppb.New(m.UpdatedAt),
	}, nil
}
