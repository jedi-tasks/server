package model

import (
	"time"

	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

// UpdateListSectionReq меняет только переданные поля
type UpdateListSectionReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID  `json:"id" validate:"required"`        // Идентификатор
	ListID    *uuid.UUID `json:"listID"`                        // Идентификатор списка
	Name      *string    `json:"name"`                          // Название раздела
	SortOrder *float64   `json:"sortOrder"`                     // Порядок раздела
	UpdatedAt time.Time  `json:"updatedAt" validate:"required"` // Время изменения на клиенте
}

// ProtoUpdateListSectionReq обертка над proto-запросом
type ProtoUpdateListSectionReq struct {
	*proto.UpdateListSectionRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoUpdateListSectionReq) ConvertToModel() (UpdateListSectionReq, error) {
	var res UpdateListSectionReq

	if p.UpdateListSectionRequest == nil {
		return res, errors.BadRequest.New("UpdateListSectionRequest is required")
	}

	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	var listID *uuid.UUID
	if p.ListID != nil {
		parsed, err := uuid.FromBytes(p.ListID)
		if err != nil {
			return res, errors.BadRequest.Wrap(err)
		}
		listID = &parsed
	}

	if p.UpdatedAt == nil {
		return res, errors.BadRequest.New("UpdatedAt is required")
	}

	return UpdateListSectionReq{
		ID:        id,
		ListID:    listID,
		Name:      p.Name,
		SortOrder: p.SortOrder,
		UpdatedAt: p.UpdatedAt.AsTime(),
	}, nil
}
