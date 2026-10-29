package model

import (
	"time"

	"server/internal/enum/sortMode"
	"server/internal/utils/errors"
	"server/internal/utils/necessary"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// UpdateSettingsReq меняет только переданные поля
type UpdateSettingsReq struct {
	Necessary            necessary.NecessaryUserInformation
	DailyCapacityMinutes *uint32            `json:"dailyCapacityMinutes"`          // Дневной лимит загрузки в минутах
	TodaySortMode        *sortMode.SortMode `json:"todaySortMode"`                 // Режим сортировки экрана «Сегодня»
	UpdatedAt            time.Time          `json:"updatedAt" validate:"required"` // Время изменения на клиенте
}

// ProtoUpdateSettingsReq обертка над proto-запросом
type ProtoUpdateSettingsReq struct {
	*proto.UpdateSettingsRequest
}

// ConvertToModel преобразует proto-запрос во внутреннюю модель
func (p ProtoUpdateSettingsReq) ConvertToModel() (UpdateSettingsReq, error) {
	var res UpdateSettingsReq

	if p.UpdateSettingsRequest == nil {
		return res, errors.BadRequest.New("UpdateSettingsRequest is required")
	}

	var todaySortMode *sortMode.SortMode
	if p.TodaySortMode != nil {
		converted, err := sortMode.ProtoSortMode{SortMode: *p.TodaySortMode}.ConvertToModel()
		if err != nil {
			return res, err
		}
		todaySortMode = &converted
	}

	if p.UpdatedAt == nil {
		return res, errors.BadRequest.New("UpdatedAt is required")
	}

	return UpdateSettingsReq{
		DailyCapacityMinutes: p.DailyCapacityMinutes,
		TodaySortMode:        todaySortMode,
		UpdatedAt:            p.UpdatedAt.AsTime(),
	}, nil
}
