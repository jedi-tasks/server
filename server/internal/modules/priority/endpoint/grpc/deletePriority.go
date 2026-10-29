package grpc

import (
	"context"

	"server/internal/modules/priority/model"
	"server/internal/utils/necessary"
	"server/internal/utils/validate"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// DeletePriority удаление записи
func (e *PriorityEndpoint) DeletePriority(ctx context.Context, r *proto.DeletePriorityRequest) (*proto.DeletePriorityResponse, error) {
	res := &proto.DeletePriorityResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoDeletePriorityReq{DeletePriorityRequest: r}.ConvertToModel()
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
	return res, e.priorityService.DeletePriority(ctx, req)
}
