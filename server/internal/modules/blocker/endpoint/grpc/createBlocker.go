package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/blocker/model"
	"server/internal/utils/necessary"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// CreateBlocker создание записи
func (e *BlockerEndpoint) CreateBlocker(ctx context.Context, r *proto.CreateBlockerRequest) (*proto.CreateBlockerResponse, error) {
	res := &proto.CreateBlockerResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoCreateBlockerReq{CreateBlockerRequest: r}.ConvertToModel()
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
	return res, e.blockerService.CreateBlocker(ctx, req)
}
