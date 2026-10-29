package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/taskList/model"
	"server/internal/utils/necessary"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// CreateTaskList создание записи
func (e *TaskListEndpoint) CreateTaskList(ctx context.Context, r *proto.CreateTaskListRequest) (*proto.CreateTaskListResponse, error) {
	res := &proto.CreateTaskListResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoCreateTaskListReq{CreateTaskListRequest: r}.ConvertToModel()
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
	return res, e.taskListService.CreateTaskList(ctx, req)
}
