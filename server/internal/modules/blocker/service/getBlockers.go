package service

import (
	"context"

	"server/internal/modules/blocker/model"
)

// GetBlockers возвращает записи пользователя по фильтрам
func (s *BlockerService) GetBlockers(ctx context.Context, req model.GetBlockersReq) ([]model.Blocker, error) {
	ctx, span := tracer.Start(ctx, "GetBlockers")
	defer span.End()

	return s.blockerRepository.GetBlockers(ctx, req)
}
