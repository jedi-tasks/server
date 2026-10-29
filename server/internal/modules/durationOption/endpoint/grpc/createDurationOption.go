package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/durationOption/model"
	"server/internal/utils/necessary"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// CreateDurationOption создание записи
func (e *DurationOptionEndpoint) CreateDurationOption(ctx context.Context, r *proto.CreateDurationOptionRequest) (*proto.CreateDurationOptionResponse, error) {
	res := &proto.CreateDurationOptionResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoCreateDurationOptionReq{CreateDurationOptionRequest: r}.ConvertToModel()
	if err != nil {
		return res, err
	}

	// Получаем необходимую информацию из контекста
	if err := necessary.ParseNecessary(ctx, &req); err != nil {
		return res, err
	}

	// Валидируем запрос
	if err := validator.Validate(req); err != nil {
		return res, err
	}

	// Вызываем метод сервиса
	return res, e.durationOptionService.CreateDurationOption(ctx, req)
}
