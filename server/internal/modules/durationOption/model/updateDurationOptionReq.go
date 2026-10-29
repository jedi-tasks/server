package model

import (
	"time"

	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

// UpdateDurationOptionReq меняет только переданные поля
type UpdateDurationOptionReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID `json:"id" validate:"required"`        // Идентификатор
	Title     *string   `json:"title"`                         // Подпись
	Minutes   *int32    `json:"minutes"`                       // Минуты, меньше нуля — непрогнозируемо долго
	UpdatedAt time.Time `json:"updatedAt" validate:"required"` // Время изменения на клиенте
}

// ProtoUpdateDurationOptionReq обертка над proto-запросом
type ProtoUpdateDurationOptionReq struct {
	*proto.UpdateDurationOptionRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoUpdateDurationOptionReq) ConvertToModel() (UpdateDurationOptionReq, error) {
	var res UpdateDurationOptionReq

	if p.UpdateDurationOptionRequest == nil {
		return res, errors.BadRequest.New("UpdateDurationOptionRequest is required")
	}

	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	if p.UpdatedAt == nil {
		return res, errors.BadRequest.New("UpdatedAt is required")
	}

	return UpdateDurationOptionReq{
		ID:        id,
		Title:     p.Title,
		Minutes:   p.Minutes,
		UpdatedAt: p.UpdatedAt.AsTime(),
	}, nil
}
