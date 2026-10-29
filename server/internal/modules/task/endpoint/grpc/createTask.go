package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/task/model"
	"server/internal/utils/necessary"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// CreateTask создание записи
func (e *TaskEndpoint) CreateTask(ctx context.Context, r *proto.CreateTaskRequest) (*proto.CreateTaskResponse, error) {
	res := &proto.CreateTaskResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoCreateTaskReq{CreateTaskRequest: r}.ConvertToModel()
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
	return res, e.taskService.CreateTask(ctx, req)
}
