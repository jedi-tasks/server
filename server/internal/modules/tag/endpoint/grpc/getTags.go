package grpc

import (
	"context"

	"server/internal/modules/tag/model"
	"server/internal/utils/necessary"
	"server/internal/utils/validate"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// GetTags получение записей
func (e *TagEndpoint) GetTags(ctx context.Context, r *proto.GetTagsRequest) (*proto.GetTagsResponse, error) {
	res := &proto.GetTagsResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoGetTagsReq{GetTagsRequest: r}.ConvertToModel()
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
	items, err := e.tagService.GetTags(ctx, req)
	if err != nil {
		return res, err
	}

	// Конвертируем ответ в proto-формат
	res.Tags = make([]*proto.Tag, 0, len(items))
	for _, item := range items {
		protoItem, err := item.ConvertToProto()
		if err != nil {
			return res, err
		}
		res.Tags = append(res.Tags, protoItem)
	}

	return res, nil
}
