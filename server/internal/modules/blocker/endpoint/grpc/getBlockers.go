package grpc

import (
	"context"

	"pkg/validator"
	"server/internal/modules/blocker/model"
	"server/internal/utils/necessary"

	proto "github.com/jedi-tasks/go-server-grpc/proto"
)

// GetBlockers получение записей
func (e *BlockerEndpoint) GetBlockers(ctx context.Context, r *proto.GetBlockersRequest) (*proto.GetBlockersResponse, error) {
	res := &proto.GetBlockersResponse{}

	// Конвертируем proto-запрос во внутреннюю модель
	req, err := model.ProtoGetBlockersReq{GetBlockersRequest: r}.ConvertToModel()
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
	items, err := e.blockerService.GetBlockers(ctx, req)
	if err != nil {
		return res, err
	}

	// Конвертируем ответ в proto-формат
	res.Blockers = make([]*proto.Blocker, 0, len(items))
	for _, item := range items {
		protoItem, err := item.ConvertToProto()
		if err != nil {
			return res, err
		}
		res.Blockers = append(res.Blockers, protoItem)
	}

	return res, nil
}
