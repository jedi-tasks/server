package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/listSection/model"
	"server/internal/utils/necessary"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// UpdateListSection редактирование записи
func (e *ListSectionEndpoint) UpdateListSection(ctx context.Context, r *proto.UpdateListSectionRequest) (*proto.UpdateListSectionResponse, error) {
	res := &proto.UpdateListSectionResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoUpdateListSectionReq{UpdateListSectionRequest: r}.ConvertToModel()
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
	return res, e.listSectionService.UpdateListSection(ctx, req)
}
