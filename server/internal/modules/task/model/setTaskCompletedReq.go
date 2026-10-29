package model

import (
	"time"

	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type SetTaskCompletedReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID `json:"id" validate:"required"`        // Идентификатор задачи
	Completed bool      `json:"completed"`                     // true — выполнить, false — вернуть в работу
	UpdatedAt time.Time `json:"updatedAt" validate:"required"` // Время действия на клиенте, оно же время выполнения
}

// ProtoSetTaskCompletedReq обертка над proto-запросом
type ProtoSetTaskCompletedReq struct {
	*proto.SetTaskCompletedRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoSetTaskCompletedReq) ConvertToModel() (SetTaskCompletedReq, error) {
	var res SetTaskCompletedReq

	if p.SetTaskCompletedRequest == nil {
		return res, errors.BadRequest.New("SetTaskCompletedRequest is required")
	}

	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	if p.UpdatedAt == nil {
		return res, errors.BadRequest.New("UpdatedAt is required")
	}

	return SetTaskCompletedReq{
		ID:        id,
		Completed: p.Completed,
		UpdatedAt: p.UpdatedAt.AsTime(),
	}, nil
}
