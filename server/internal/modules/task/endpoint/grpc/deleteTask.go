package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/task/model"
	"server/internal/utils/necessary"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// DeleteTask удаление записи
func (e *TaskEndpoint) DeleteTask(ctx context.Context, r *proto.DeleteTaskRequest) (*proto.DeleteTaskResponse, error) {
	res := &proto.DeleteTaskResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoDeleteTaskReq{DeleteTaskRequest: r}.ConvertToModel()
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
	return res, e.taskService.DeleteTask(ctx, req)
}
