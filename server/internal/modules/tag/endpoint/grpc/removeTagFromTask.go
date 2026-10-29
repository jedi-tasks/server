package grpc

import (
	"context"

	"server/internal/modules/tag/model"
	"server/internal/utils/necessary"
	"server/internal/utils/validate"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// RemoveTagFromTask привязка/отвязка тега
func (e *TagEndpoint) RemoveTagFromTask(ctx context.Context, r *proto.RemoveTagFromTaskRequest) (*proto.RemoveTagFromTaskResponse, error) {
	res := &proto.RemoveTagFromTaskResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoRemoveTagFromTaskReq{RemoveTagFromTaskRequest: r}.ConvertToModel()
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
	return res, e.tagService.RemoveTagFromTask(ctx, req)
}
