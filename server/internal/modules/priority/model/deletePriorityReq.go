package model

import (
	"time"

	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type DeletePriorityReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID `json:"id" validate:"required"`        // Идентификатор
	UpdatedAt time.Time `json:"updatedAt" validate:"required"` // Время удаления на клиенте
}

// ProtoDeletePriorityReq обертка над proto-запросом
type ProtoDeletePriorityReq struct {
	*proto.DeletePriorityRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoDeletePriorityReq) ConvertToModel() (DeletePriorityReq, error) {
	var res DeletePriorityReq

	if p.DeletePriorityRequest == nil {
		return res, errors.BadRequest.New("DeletePriorityRequest is required")
	}

	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	if p.UpdatedAt == nil {
		return res, errors.BadRequest.New("UpdatedAt is required")
	}

	return DeletePriorityReq{
		ID:        id,
		UpdatedAt: p.UpdatedAt.AsTime(),
	}, nil
}
