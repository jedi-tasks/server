package grpc

import (
	"context"

	"server/internal/modules/priority/model"
	"server/internal/utils/necessary"
	"server/internal/utils/validate"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// GetPriorities получение записей
func (e *PriorityEndpoint) GetPriorities(ctx context.Context, r *proto.GetPrioritiesRequest) (*proto.GetPrioritiesResponse, error) {
	res := &proto.GetPrioritiesResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoGetPrioritiesReq{GetPrioritiesRequest: r}.ConvertToModel()
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
	items, err := e.priorityService.GetPriorities(ctx, req)
	if err != nil {
		return res, err
	}

	// Конвертируем ответ в proto-формат
	res.Priorities = make([]*proto.Priority, 0, len(items))
	for _, item := range items {
		protoItem, err := item.ConvertToProto()
		if err != nil {
			return res, err
		}
		res.Priorities = append(res.Priorities, protoItem)
	}

	return res, nil
}
