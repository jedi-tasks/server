package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/task/model"
	"server/internal/utils/necessary"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// UpdateTask редактирование записи
func (e *TaskEndpoint) UpdateTask(ctx context.Context, r *proto.UpdateTaskRequest) (*proto.UpdateTaskResponse, error) {
	res := &proto.UpdateTaskResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoUpdateTaskReq{UpdateTaskRequest: r}.ConvertToModel()
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
	return res, e.taskService.UpdateTask(ctx, req)
}
