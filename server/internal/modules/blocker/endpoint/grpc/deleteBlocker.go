package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/blocker/model"
	"server/internal/utils/necessary"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// DeleteBlocker удаление записи
func (e *BlockerEndpoint) DeleteBlocker(ctx context.Context, r *proto.DeleteBlockerRequest) (*proto.DeleteBlockerResponse, error) {
	res := &proto.DeleteBlockerResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoDeleteBlockerReq{DeleteBlockerRequest: r}.ConvertToModel()
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
	return res, e.blockerService.DeleteBlocker(ctx, req)
}
