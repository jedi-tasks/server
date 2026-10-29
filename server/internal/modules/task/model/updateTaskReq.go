package model

import (
	"context"
	"time"

	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

// UpdateTaskReq меняет только переданные поля
type UpdateTaskReq struct {
	Necessary       necessary.NecessaryUserInformation
	ID              uuid.UUID  `json:"id" validate:"required"`        // Идентификатор
	ListID          *uuid.UUID `json:"listID"`                        // Идентификатор списка
	SectionID       *uuid.UUID `json:"sectionID"`                     // Идентификатор раздела uuid.Nil — сбросить
	Title           *string    `json:"title"`                         // Название задачи
	Notes           *string    `json:"notes"`                         // Описание
	PriorityID      *uuid.UUID `json:"priorityID"`                    // Идентификатор приоритета uuid.Nil — сбросить
	PlannedDay      *string    `json:"plannedDay"`                    // Планируемый срок, yyyy-MM-dd пустая строка — сбросить
	DeadlineDay     *string    `json:"deadlineDay"`                   // Дедлайн, yyyy-MM-dd пустая строка — сбросить
	DurationMinutes *int32     `json:"durationMinutes"`               // Планируемое время в минутах 0 — сбросить
	ManualOrder     *float64   `json:"manualOrder"`                   // Ручной порядок
	UpdatedAt       time.Time  `json:"updatedAt" validate:"required"` // Время изменения на клиенте
}

// Validate проверяет корректность запроса
func (s UpdateTaskReq) Validate(ctx context.Context) error {
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

// ProtoUpdateTaskReq обертка над proto-запросом
type ProtoUpdateTaskReq struct {
	*proto.UpdateTaskRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoUpdateTaskReq) ConvertToModel() (UpdateTaskReq, error) {
	var res UpdateTaskReq

	if p.UpdateTaskRequest == nil {
		return res, errors.BadRequest.New("UpdateTaskRequest is required")
	}

	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	var listID *uuid.UUID
	if p.ListID != nil {
		parsed, err := uuid.FromBytes(p.ListID)
		if err != nil {
			return res, errors.BadRequest.Wrap(err)
		}
		listID = &parsed
	}
	var sectionID *uuid.UUID
	if p.SectionID != nil {
		if len(p.SectionID) == 0 {
			sectionID = &uuid.Nil
		} else {
			parsed, err := uuid.FromBytes(p.SectionID)
			if err != nil {
				return res, errors.BadRequest.Wrap(err)
			}
			sectionID = &parsed
		}
	}
	var priorityID *uuid.UUID
	if p.PriorityID != nil {
		if len(p.PriorityID) == 0 {
			priorityID = &uuid.Nil
		} else {
			parsed, err := uuid.FromBytes(p.PriorityID)
			if err != nil {
				return res, errors.BadRequest.Wrap(err)
			}
			priorityID = &parsed
		}
	}

	if p.UpdatedAt == nil {
		return res, errors.BadRequest.New("UpdatedAt is required")
	}

	return UpdateTaskReq{
		ID:              id,
		ListID:          listID,
		SectionID:       sectionID,
		Title:           p.Title,
		Notes:           p.Notes,
		PriorityID:      priorityID,
		PlannedDay:      p.PlannedDay,
		DeadlineDay:     p.DeadlineDay,
		DurationMinutes: p.DurationMinutes,
		ManualOrder:     p.ManualOrder,
		UpdatedAt:       p.UpdatedAt.AsTime(),
	}, nil
}
