package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/listSection/model"
	"server/internal/utils/necessary"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// GetListSections получение записей
func (e *ListSectionEndpoint) GetListSections(ctx context.Context, r *proto.GetListSectionsRequest) (*proto.GetListSectionsResponse, error) {
	res := &proto.GetListSectionsResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoGetListSectionsReq{GetListSectionsRequest: r}.ConvertToModel()
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
	items, err := e.listSectionService.GetListSections(ctx, req)
	if err != nil {
		return res, err
	}

	// Конвертируем ответ в proto-формат
	res.ListSections = make([]*proto.ListSection, 0, len(items))
	for _, item := range items {
		protoItem, err := item.ConvertToProto()
		if err != nil {
			return res, err
		}
		res.ListSections = append(res.ListSections, protoItem)
	}

	return res, nil
}
