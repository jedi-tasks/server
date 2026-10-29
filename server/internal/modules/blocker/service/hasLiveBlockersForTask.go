package service

import (
	"context"

	"github.com/google/uuid"
)

// HasLiveBlockersForTask проверяет, есть ли неудаленные блокеры, где задача заблокирована или блокирует
func (s *BlockerService) HasLiveBlockersForTask(ctx context.Context, taskID uuid.UUID) (bool, error) {
	ctx, span := tracer.Start(ctx, "HasLiveBlockersForTask")
	defer span.End()

	return s.blockerRepository.HasLiveBlockersForTask(ctx, taskID)
}
