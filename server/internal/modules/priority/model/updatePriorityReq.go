package model

import (
	"time"

	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

// UpdatePriorityReq меняет только переданные поля
type UpdatePriorityReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID `json:"id" validate:"required"`        // Идентификатор
	Name      *string   `json:"name"`                          // Название
	ColorHex  *string   `json:"colorHex"`                      // Цвет
	Marks     *string   `json:"marks"`                         // Знаки в строке задачи
	SortOrder *float64  `json:"sortOrder"`                     // Порядок важности
	UpdatedAt time.Time `json:"updatedAt" validate:"required"` // Время изменения на клиенте
}

// ProtoUpdatePriorityReq обертка над proto-запросом
type ProtoUpdatePriorityReq struct {
	*proto.UpdatePriorityRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoUpdatePriorityReq) ConvertToModel() (UpdatePriorityReq, error) {
	var res UpdatePriorityReq

	if p.UpdatePriorityRequest == nil {
		return res, errors.BadRequest.New("UpdatePriorityRequest is required")
	}

	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	if p.UpdatedAt == nil {
		return res, errors.BadRequest.New("UpdatedAt is required")
	}

	return UpdatePriorityReq{
		ID:        id,
		Name:      p.Name,
		ColorHex:  p.ColorHex,
		Marks:     p.Marks,
		SortOrder: p.SortOrder,
		UpdatedAt: p.UpdatedAt.AsTime(),
	}, nil
}
