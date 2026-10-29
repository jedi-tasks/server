package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/tag/model"
	"server/internal/utils/necessary"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// GetTasksToTags получение связей задач и тегов
func (e *TagEndpoint) GetTasksToTags(ctx context.Context, r *proto.GetTasksToTagsRequest) (*proto.GetTasksToTagsResponse, error) {
	res := &proto.GetTasksToTagsResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoGetTasksToTagsReq{GetTasksToTagsRequest: r}.ConvertToModel()
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
	links, err := e.tagService.GetTasksToTags(ctx, req)
	if err != nil {
		return res, err
	}

	// Конвертируем ответ в proto-формат
	res.TasksToTags = make([]*proto.TaskToTag, 0, len(links))
	for _, link := range links {
		res.TasksToTags = append(res.TasksToTags, link.ConvertToProto())
	}

	return res, nil
}
