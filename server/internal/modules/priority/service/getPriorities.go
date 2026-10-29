package service

import (
	"context"

	"server/internal/modules/priority/model"
)

// GetPriorities возвращает записи пользователя по фильтрам
func (s *PriorityService) GetPriorities(ctx context.Context, req model.GetPrioritiesReq) ([]model.Priority, error) {
	ctx, span := tracer.Start(ctx, "GetPriorities")
	defer span.End()

	return s.priorityRepository.GetPriorities(ctx, req)
}
