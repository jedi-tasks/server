package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/tag/model"
	"server/internal/utils/necessary"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// DeleteTag удаление записи
func (e *TagEndpoint) DeleteTag(ctx context.Context, r *proto.DeleteTagRequest) (*proto.DeleteTagResponse, error) {
	res := &proto.DeleteTagResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoDeleteTagReq{DeleteTagRequest: r}.ConvertToModel()
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
	return res, e.tagService.DeleteTag(ctx, req)
}
