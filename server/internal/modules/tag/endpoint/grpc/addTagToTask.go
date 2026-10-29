package grpc

import (
	"context"

	"server/internal/modules/tag/model"
	"server/internal/utils/necessary"
	"server/internal/utils/validate"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// AddTagToTask привязка/отвязка тега
func (e *TagEndpoint) AddTagToTask(ctx context.Context, r *proto.AddTagToTaskRequest) (*proto.AddTagToTaskResponse, error) {
	res := &proto.AddTagToTaskResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoAddTagToTaskReq{AddTagToTaskRequest: r}.ConvertToModel()
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
	return res, e.tagService.AddTagToTask(ctx, req)
}
