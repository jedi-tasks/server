package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/taskList/model"
	"server/internal/utils/necessary"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// UpdateTaskList редактирование записи
func (e *TaskListEndpoint) UpdateTaskList(ctx context.Context, r *proto.UpdateTaskListRequest) (*proto.UpdateTaskListResponse, error) {
	res := &proto.UpdateTaskListResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoUpdateTaskListReq{UpdateTaskListRequest: r}.ConvertToModel()
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
	return res, e.taskListService.UpdateTaskList(ctx, req)
}
