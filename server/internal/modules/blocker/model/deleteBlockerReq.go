package model

import (
	"time"

	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type DeleteBlockerReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID `json:"id" validate:"required"`        // Идентификатор
	UpdatedAt time.Time `json:"updatedAt" validate:"required"` // Время удаления на клиенте
}

// ProtoDeleteBlockerReq обертка над proto-запросом
type ProtoDeleteBlockerReq struct {
	*proto.DeleteBlockerRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoDeleteBlockerReq) ConvertToModel() (DeleteBlockerReq, error) {
	var res DeleteBlockerReq

	if p.DeleteBlockerRequest == nil {
		return res, errors.BadRequest.New("DeleteBlockerRequest is required")
	}

	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	if p.UpdatedAt == nil {
		return res, errors.BadRequest.New("UpdatedAt is required")
	}

	return DeleteBlockerReq{
		ID:        id,
		UpdatedAt: p.UpdatedAt.AsTime(),
	}, nil
}
