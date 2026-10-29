package model

import (
	"time"

	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type DeleteTaskListReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID `json:"id" validate:"required"`        // Идентификатор
	UpdatedAt time.Time `json:"updatedAt" validate:"required"` // Время удаления на клиенте
}

// ProtoDeleteTaskListReq обертка над proto-запросом
type ProtoDeleteTaskListReq struct {
	*proto.DeleteTaskListRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoDeleteTaskListReq) ConvertToModel() (DeleteTaskListReq, error) {
	var res DeleteTaskListReq

	if p.DeleteTaskListRequest == nil {
		return res, errors.BadRequest.New("DeleteTaskListRequest is required")
	}

	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	if p.UpdatedAt == nil {
		return res, errors.BadRequest.New("UpdatedAt is required")
	}

	return DeleteTaskListReq{
		ID:        id,
		UpdatedAt: p.UpdatedAt.AsTime(),
	}, nil
}
