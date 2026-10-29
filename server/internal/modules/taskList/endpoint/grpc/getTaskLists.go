package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/taskList/model"
	"server/internal/utils/necessary"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// GetTaskLists получение записей
func (e *TaskListEndpoint) GetTaskLists(ctx context.Context, r *proto.GetTaskListsRequest) (*proto.GetTaskListsResponse, error) {
	res := &proto.GetTaskListsResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoGetTaskListsReq{GetTaskListsRequest: r}.ConvertToModel()
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
	items, err := e.taskListService.GetTaskLists(ctx, req)
	if err != nil {
		return res, err
	}

	// Конвертируем ответ в proto-формат
	res.TaskLists = make([]*proto.TaskList, 0, len(items))
	for _, item := range items {
		protoItem, err := item.ConvertToProto()
		if err != nil {
			return res, err
		}
		res.TaskLists = append(res.TaskLists, protoItem)
	}

	return res, nil
}
