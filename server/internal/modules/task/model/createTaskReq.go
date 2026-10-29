package model

import (
	"context"
	"time"

	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	repoModel "server/internal/modules/task/repository/model"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type CreateTaskReq struct {
	Necessary       necessary.NecessaryUserInformation
	ID              uuid.UUID  `json:"id" validate:"required"`        // Идентификатор, генерируется клиентом
	ListID          uuid.UUID  `json:"listID"`                        // Идентификатор списка, uuid.Nil — Inbox
	SectionID       *uuid.UUID `json:"sectionID"`                     // Идентификатор раздела
	Title           string     `json:"title"`                         // Название задачи
	Notes           string     `json:"notes"`                         // Описание
	PriorityID      *uuid.UUID `json:"priorityID"`                    // Идентификатор приоритета
	PlannedDay      *string    `json:"plannedDay"`                    // Планируемый срок, yyyy-MM-dd
	DeadlineDay     *string    `json:"deadlineDay"`                   // Дедлайн, yyyy-MM-dd
	DurationMinutes *int32     `json:"durationMinutes"`               // Планируемое время в минутах
	ManualOrder     float64    `json:"manualOrder"`                   // Ручной порядок
	ExternalID      *string    `json:"externalID"`                    // Идентификатор напоминания при импорте
	CreatedAt       time.Time  `json:"createdAt" validate:"required"` // Время создания на клиенте
}

// Validate проверяет корректность запроса
func (s CreateTaskReq) Validate(ctx context.Context) error {
	if s.PlannedDay != nil && *s.PlannedDay != "" {
		if _, err := time.Parse(time.DateOnly, *s.PlannedDay); err != nil {
			return errors.BadRequest.New("plannedDay must be yyyy-MM-dd").WithContextParams(ctx)
		}
	}
	if s.DeadlineDay != nil && *s.DeadlineDay != "" {
		if _, err := time.Parse(time.DateOnly, *s.DeadlineDay); err != nil {
			return errors.BadRequest.New("deadlineDay must be yyyy-MM-dd").WithContextParams(ctx)
		}
	}
	return nil
}

// ConvertToRepoReq преобразует запрос в модель репозитория
func (s CreateTaskReq) ConvertToRepoReq() repoModel.CreateTaskReq {
	return repoModel.CreateTaskReq{
		ID:              s.ID,
		UserID:          s.Necessary.UserID,
		ListID:          s.ListID,
		SectionID:       s.SectionID,
		Title:           s.Title,
		Notes:           s.Notes,
		PriorityID:      s.PriorityID,
		PlannedDay:      s.PlannedDay,
		DeadlineDay:     s.DeadlineDay,
		DurationMinutes: s.DurationMinutes,
		ManualOrder:     s.ManualOrder,
		ExternalID:      s.ExternalID,
		CreatedAt:       s.CreatedAt,
	}
}

// ProtoCreateTaskReq обертка над proto-запросом
type ProtoCreateTaskReq struct {
	*proto.CreateTaskRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoCreateTaskReq) ConvertToModel() (CreateTaskReq, error) {
	var res CreateTaskReq

	if p.CreateTaskRequest == nil {
		return res, errors.BadRequest.New("CreateTaskRequest is required")
	}

	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	// Без списка задача попадает в Inbox
	listID := uuid.Nil
	if len(p.ListID) != 0 {
		listID, err = uuid.FromBytes(p.ListID)
		if err != nil {
			return res, errors.BadRequest.Wrap(err)
		}
	}
	var sectionID *uuid.UUID
	if len(p.SectionID) != 0 {
		parsed, err := uuid.FromBytes(p.SectionID)
		if err != nil {
			return res, errors.BadRequest.Wrap(err)
		}
		sectionID = &parsed
	}
	var priorityID *uuid.UUID
	if len(p.PriorityID) != 0 {
		parsed, err := uuid.FromBytes(p.PriorityID)
		if err != nil {
			return res, errors.BadRequest.Wrap(err)
		}
		priorityID = &parsed
	}
	var plannedDay *string
	if p.PlannedDay != nil && *p.PlannedDay != "" {
		plannedDay = p.PlannedDay
	}
	var deadlineDay *string
	if p.DeadlineDay != nil && *p.DeadlineDay != "" {
		deadlineDay = p.DeadlineDay
	}
	var durationMinutes *int32
	if p.DurationMinutes != nil && *p.DurationMinutes != 0 {
		durationMinutes = p.DurationMinutes
	}
	var externalID *string
	if p.ExternalID != nil && *p.ExternalID != "" {
		externalID = p.ExternalID
	}
	if p.CreatedAt == nil {
		return res, errors.BadRequest.New("CreatedAt is required")
	}

	return CreateTaskReq{
		ID:              id,
		ListID:          listID,
		SectionID:       sectionID,
		Title:           p.Title,
		Notes:           p.Notes,
		PriorityID:      priorityID,
		PlannedDay:      plannedDay,
		DeadlineDay:     deadlineDay,
		DurationMinutes: durationMinutes,
		ManualOrder:     p.ManualOrder,
		ExternalID:      externalID,
		CreatedAt:       p.CreatedAt.AsTime(),
	}, nil
}
