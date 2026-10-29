package model

import (
	"time"

	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

// UpdateTagReq меняет только переданные поля
type UpdateTagReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID `json:"id" validate:"required"`        // Идентификатор
	Name      *string   `json:"name"`                          // Название тега
	UpdatedAt time.Time `json:"updatedAt" validate:"required"` // Время изменения на клиенте
}

// ProtoUpdateTagReq обертка над proto-запросом
type ProtoUpdateTagReq struct {
	*proto.UpdateTagRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoUpdateTagReq) ConvertToModel() (UpdateTagReq, error) {
	var res UpdateTagReq

	if p.UpdateTagRequest == nil {
		return res, errors.BadRequest.New("UpdateTagRequest is required")
	}

	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	if p.UpdatedAt == nil {
		return res, errors.BadRequest.New("UpdatedAt is required")
	}

	return UpdateTagReq{
		ID:        id,
		Name:      p.Name,
		UpdatedAt: p.UpdatedAt.AsTime(),
	}, nil
}
