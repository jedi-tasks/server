package model

import (
	"time"

	"github.com/google/uuid"

	proto "github.com/jedi-tasks/go-server-grpc/proto"

	"server/internal/enum/sortMode"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"
)

// UpdateTaskListReq меняет только переданные поля
type UpdateTaskListReq struct {
	Necessary necessary.NecessaryUserInformation
	ID        uuid.UUID          `json:"id" validate:"required"`        // Идентификатор
	Name      *string            `json:"name"`                          // Название списка
	ColorHex  *string            `json:"colorHex"`                      // Цвет списка
	Icon      *string            `json:"icon"`                          // SF Symbol иконки списка
	SortOrder *float64           `json:"sortOrder"`                     // Порядок списка
	SortMode  *sortMode.SortMode `json:"sortMode"`                      // Режим сортировки задач в списке
	UpdatedAt time.Time          `json:"updatedAt" validate:"required"` // Время изменения на клиенте
}

// ProtoUpdateTaskListReq обертка над proto-запросом
type ProtoUpdateTaskListReq struct {
	*proto.UpdateTaskListRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoUpdateTaskListReq) ConvertToModel() (UpdateTaskListReq, error) {
	var res UpdateTaskListReq

	if p.UpdateTaskListRequest == nil {
		return res, errors.BadRequest.New("UpdateTaskListRequest is required")
	}

	id, err := uuid.FromBytes(p.Id)
	if err != nil {
		return res, errors.BadRequest.Wrap(err)
	}

	var sortModeValue *sortMode.SortMode
	if p.SortMode != nil {
		converted, err := sortMode.ProtoSortMode{SortMode: *p.SortMode}.ConvertToModel()
		if err != nil {
			return res, err
		}
		sortModeValue = &converted
	}

	if p.UpdatedAt == nil {
		return res, errors.BadRequest.New("UpdatedAt is required")
	}

	return UpdateTaskListReq{
		ID:        id,
		Name:      p.Name,
		ColorHex:  p.ColorHex,
		Icon:      p.Icon,
		SortOrder: p.SortOrder,
		SortMode:  sortModeValue,
		UpdatedAt: p.UpdatedAt.AsTime(),
	}, nil
}
