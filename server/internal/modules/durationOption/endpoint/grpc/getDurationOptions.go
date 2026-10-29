package grpc

import (
	"context"

	"server/internal/modules/durationOption/model"
	"server/internal/utils/necessary"
	"server/internal/utils/validate"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// GetDurationOptions получение записей
func (e *DurationOptionEndpoint) GetDurationOptions(ctx context.Context, r *proto.GetDurationOptionsRequest) (*proto.GetDurationOptionsResponse, error) {
	res := &proto.GetDurationOptionsResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoGetDurationOptionsReq{GetDurationOptionsRequest: r}.ConvertToModel()
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
	items, err := e.durationOptionService.GetDurationOptions(ctx, req)
	if err != nil {
		return res, err
	}

	// Конвертируем ответ в proto-формат
	res.DurationOptions = make([]*proto.DurationOption, 0, len(items))
	for _, item := range items {
		protoItem, err := item.ConvertToProto()
		if err != nil {
			return res, err
		}
		res.DurationOptions = append(res.DurationOptions, protoItem)
	}

	return res, nil
}
