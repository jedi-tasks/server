package model

import (
	"context"
	"time"

	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	repoModel "server/internal/modules/blocker/repository/model"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type CreateBlockerReq struct {
	Necessary      necessary.NecessaryUserInformation
	ID             uuid.UUID  `json:"id" validate:"required"`        // Идентификатор, генерируется клиентом
	TaskID         uuid.UUID  `json:"taskID" validate:"required"`    // Заблокированная задача
	BlockingTaskID *uuid.UUID `json:"blockingTaskID"`                // Блокирующая задача
	Reason         *string    `json:"reason"`                        // Текстовая причина
	SinceDay       string     `json:"sinceDay"`                      // Дата начала, yyyy-MM-dd
	UpdatedAt      time.Time  `json:"updatedAt" validate:"required"` // Время создания на клиенте
}

// Validate проверяет корректность запроса
func (s CreateBlockerReq) Validate(ctx context.Context) error {
	if (s.BlockingTaskID == nil) == (s.Reason == nil) {
		return errors.BadRequest.New("exactly one of blockingTaskID and reason must be set").WithContextParams(ctx)
	}
	if s.BlockingTaskID != nil && *s.BlockingTaskID == s.TaskID {
		return errors.BadRequest.New("task can't block itself").WithContextParams(ctx)
	}
	if _, err := time.Parse(time.DateOnly, s.SinceDay); err != nil {
		return errors.BadRequest.New("sinceDay must be yyyy-MM-dd").WithContextParams(ctx)
	}
	return nil
}

// ConvertToRepoReq преобразует запрос в модель репозитория
func (s CreateBlockerReq) ConvertToRepoReq() repoModel.CreateBlockerReq {
	return repoModel.CreateBlockerReq{
		ID:             s.ID,
		UserID:         s.Necessary.UserID,
		TaskID:         s.TaskID,
		BlockingTaskID: s.BlockingTaskID,
		Reason:         s.Reason,
		SinceDay:       s.SinceDay,
		UpdatedAt:      s.UpdatedAt,
	}
}

// ProtoCreateBlockerReq обертка над proto-запросом
type ProtoCreateBlockerReq struct {
	*proto.CreateBlockerRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoCreateBlockerReq) ConvertToModel() (CreateBlockerReq, error) {
	var res CreateBlockerReq

	if p.CreateBlockerRequest == nil {
		return res, errors.BadRequest.New("CreateBlockerRequest is required")
	}

	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	taskID, err := uuid.FromBytes(p.TaskID)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}
	var blockingTaskID *uuid.UUID
	if len(p.BlockingTaskID) != 0 {
		parsed, err := uuid.FromBytes(p.BlockingTaskID)
		if err != nil {
			return res, errors.BadRequest.Wrap(err)
		}
		blockingTaskID = &parsed
	}
	var reason *string
	if p.Reason != nil && *p.Reason != "" {
		reason = p.Reason
	}
	if p.UpdatedAt == nil {
		return res, errors.BadRequest.New("UpdatedAt is required")
	}

	return CreateBlockerReq{
		ID:             id,
		TaskID:         taskID,
		BlockingTaskID: blockingTaskID,
		Reason:         reason,
		SinceDay:       p.SinceDay,
		UpdatedAt:      p.UpdatedAt.AsTime(),
	}, nil
}
