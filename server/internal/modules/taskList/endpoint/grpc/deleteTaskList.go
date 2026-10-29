package grpc

import (
	"context"

	"server/internal/modules/taskList/model"
	"server/internal/utils/necessary"
	"server/internal/utils/validate"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// DeleteTaskList удаление записи
func (e *TaskListEndpoint) DeleteTaskList(ctx context.Context, r *proto.DeleteTaskListRequest) (*proto.DeleteTaskListResponse, error) {
	res := &proto.DeleteTaskListResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoDeleteTaskListReq{DeleteTaskListRequest: r}.ConvertToModel()
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
	return res, e.taskListService.DeleteTaskList(ctx, req)
}
