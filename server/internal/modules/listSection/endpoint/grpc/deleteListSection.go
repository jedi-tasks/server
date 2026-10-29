package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/listSection/model"
	"server/internal/utils/necessary"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// DeleteListSection удаление записи
func (e *ListSectionEndpoint) DeleteListSection(ctx context.Context, r *proto.DeleteListSectionRequest) (*proto.DeleteListSectionResponse, error) {
	res := &proto.DeleteListSectionResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoDeleteListSectionReq{DeleteListSectionRequest: r}.ConvertToModel()
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
	return res, e.listSectionService.DeleteListSection(ctx, req)
}
