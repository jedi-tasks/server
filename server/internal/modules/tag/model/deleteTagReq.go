package model

import (
	"time"

	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type DeleteTagReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID `json:"id" validate:"required"`        // Идентификатор
	UpdatedAt time.Time `json:"updatedAt" validate:"required"` // Время удаления на клиенте
}

// ProtoDeleteTagReq обертка над proto-запросом
type ProtoDeleteTagReq struct {
	*proto.DeleteTagRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoDeleteTagReq) ConvertToModel() (DeleteTagReq, error) {
	var res DeleteTagReq

	if p.DeleteTagRequest == nil {
		return res, errors.BadRequest.New("DeleteTagRequest is required")
	}

	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	if p.UpdatedAt == nil {
		return res, errors.BadRequest.New("UpdatedAt is required")
	}

	return DeleteTagReq{
		ID:        id,
		UpdatedAt: p.UpdatedAt.AsTime(),
	}, nil
}
