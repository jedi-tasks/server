package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/tag/model"
	"server/internal/utils/necessary"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// CreateTag создание записи
func (e *TagEndpoint) CreateTag(ctx context.Context, r *proto.CreateTagRequest) (*proto.CreateTagResponse, error) {
	res := &proto.CreateTagResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoCreateTagReq{CreateTagRequest: r}.ConvertToModel()
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
	return res, e.tagService.CreateTag(ctx, req)
}
