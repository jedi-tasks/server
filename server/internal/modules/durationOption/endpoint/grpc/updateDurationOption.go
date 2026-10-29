package grpc

import (
	"context"

	"server/internal/modules/durationOption/model"
	"server/internal/utils/necessary"
	"server/internal/utils/validate"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// UpdateDurationOption редактирование записи
func (e *DurationOptionEndpoint) UpdateDurationOption(ctx context.Context, r *proto.UpdateDurationOptionRequest) (*proto.UpdateDurationOptionResponse, error) {
	res := &proto.UpdateDurationOptionResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoUpdateDurationOptionReq{UpdateDurationOptionRequest: r}.ConvertToModel()
	if err != nil {
		return res, err
	}

	// Получаем необходимую информацию из контекста
	if err := necessary.ParseNecessary(ctx, &req); err != nil {
		return res, err
	}

	// Валидируем запрос
	if err := validate.Validate(ctx, req); err != nil {
		return res, err
	}

	// Вызываем метод сервиса
	return res, e.durationOptionService.UpdateDurationOption(ctx, req)
}
