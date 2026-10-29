package grpc

import (
	"context"

	"server/internal/modules/task/model"
	"server/internal/utils/necessary"
	"server/internal/utils/validate"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// GetTasks получение записей
func (e *TaskEndpoint) GetTasks(ctx context.Context, r *proto.GetTasksRequest) (*proto.GetTasksResponse, error) {
	res := &proto.GetTasksResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoGetTasksReq{GetTasksRequest: r}.ConvertToModel()
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
	items, err := e.taskService.GetTasks(ctx, req)
	if err != nil {
		return res, err
	}

	// Конвертируем ответ в proto-формат
	res.Tasks = make([]*proto.Task, 0, len(items))
	for _, item := range items {
		protoItem, err := item.ConvertToProto()
		if err != nil {
			return res, err
		}
		res.Tasks = append(res.Tasks, protoItem)
	}

	return res, nil
}
