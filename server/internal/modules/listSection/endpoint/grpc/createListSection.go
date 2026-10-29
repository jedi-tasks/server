package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/listSection/model"
	"server/internal/utils/necessary"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// CreateListSection создание записи
func (e *ListSectionEndpoint) CreateListSection(ctx context.Context, r *proto.CreateListSectionRequest) (*proto.CreateListSectionResponse, error) {
	res := &proto.CreateListSectionResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoCreateListSectionReq{CreateListSectionRequest: r}.ConvertToModel()
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
	return res, e.listSectionService.CreateListSection(ctx, req)
}
