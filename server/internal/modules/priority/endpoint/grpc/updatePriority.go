package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/priority/model"
	"server/internal/utils/necessary"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// UpdatePriority редактирование записи
func (e *PriorityEndpoint) UpdatePriority(ctx context.Context, r *proto.UpdatePriorityRequest) (*proto.UpdatePriorityResponse, error) {
	res := &proto.UpdatePriorityResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoUpdatePriorityReq{UpdatePriorityRequest: r}.ConvertToModel()
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
	return res, e.priorityService.UpdatePriority(ctx, req)
}
