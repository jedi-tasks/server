package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/tag/model"
	"server/internal/utils/necessary"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// UpdateTag редактирование записи
func (e *TagEndpoint) UpdateTag(ctx context.Context, r *proto.UpdateTagRequest) (*proto.UpdateTagResponse, error) {
	res := &proto.UpdateTagResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoUpdateTagReq{UpdateTagRequest: r}.ConvertToModel()
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
	return res, e.tagService.UpdateTag(ctx, req)
}
