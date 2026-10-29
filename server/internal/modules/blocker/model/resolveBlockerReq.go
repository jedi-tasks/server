package model

import (
	"time"

	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

type ResolveBlockerReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID `json:"id" validate:"required"`        // Идентификатор блокера
	UpdatedAt time.Time `json:"updatedAt" validate:"required"` // Время снятия на клиенте
}

// ProtoResolveBlockerReq обертка над proto-запросом
type ProtoResolveBlockerReq struct {
	*proto.ResolveBlockerRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoResolveBlockerReq) ConvertToModel() (ResolveBlockerReq, error) {
	var res ResolveBlockerReq

	if p.ResolveBlockerRequest == nil {
		return res, errors.BadRequest.New("ResolveBlockerRequest is required")
	}

	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	if p.UpdatedAt == nil {
		return res, errors.BadRequest.New("UpdatedAt is required")
	}

	return ResolveBlockerReq{
		ID:        id,
		UpdatedAt: p.UpdatedAt.AsTime(),
	}, nil
}
