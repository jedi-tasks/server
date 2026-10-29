package model

import (
	"time"

	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type DeleteListSectionReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID `json:"id" validate:"required"`        // Идентификатор
	UpdatedAt time.Time `json:"updatedAt" validate:"required"` // Время удаления на клиенте
}

// ProtoDeleteListSectionReq обертка над proto-запросом
type ProtoDeleteListSectionReq struct {
	*proto.DeleteListSectionRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoDeleteListSectionReq) ConvertToModel() (DeleteListSectionReq, error) {
	var res DeleteListSectionReq

	if p.DeleteListSectionRequest == nil {
		return res, errors.BadRequest.New("DeleteListSectionRequest is required")
	}

	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	if p.UpdatedAt == nil {
		return res, errors.BadRequest.New("UpdatedAt is required")
	}

	return DeleteListSectionReq{
		ID:        id,
		UpdatedAt: p.UpdatedAt.AsTime(),
	}, nil
}
