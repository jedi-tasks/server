package model

import (
	"time"

	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type DeleteTaskReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID `json:"id" validate:"required"`        // Идентификатор
	UpdatedAt time.Time `json:"updatedAt" validate:"required"` // Время удаления на клиенте
}

// ProtoDeleteTaskReq обертка над proto-запросом
type ProtoDeleteTaskReq struct {
	*proto.DeleteTaskRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoDeleteTaskReq) ConvertToModel() (DeleteTaskReq, error) {
	var res DeleteTaskReq

	if p.DeleteTaskRequest == nil {
		return res, errors.BadRequest.New("DeleteTaskRequest is required")
	}

	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	if p.UpdatedAt == nil {
		return res, errors.BadRequest.New("UpdatedAt is required")
	}

	return DeleteTaskReq{
		ID:        id,
		UpdatedAt: p.UpdatedAt.AsTime(),
	}, nil
}
