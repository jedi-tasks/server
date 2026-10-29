package grpc

import (
	"context"

	"server/internal/modules/blocker/model"
	"server/internal/utils/necessary"
	"server/internal/utils/validate"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// ResolveBlocker снятие блокера
func (e *BlockerEndpoint) ResolveBlocker(ctx context.Context, r *proto.ResolveBlockerRequest) (*proto.ResolveBlockerResponse, error) {
	res := &proto.ResolveBlockerResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoResolveBlockerReq{ResolveBlockerRequest: r}.ConvertToModel()
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
	return res, e.blockerService.ResolveBlocker(ctx, req)
}
