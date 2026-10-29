package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/priority/model"
	"server/internal/utils/necessary"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// CreatePriority создание записи
func (e *PriorityEndpoint) CreatePriority(ctx context.Context, r *proto.CreatePriorityRequest) (*proto.CreatePriorityResponse, error) {
	res := &proto.CreatePriorityResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoCreatePriorityReq{CreatePriorityRequest: r}.ConvertToModel()
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
	return res, e.priorityService.CreatePriority(ctx, req)
}
