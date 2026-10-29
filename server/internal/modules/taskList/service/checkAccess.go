package service

import (
	"context"

	"github.com/google/uuid"
)

// CheckAccess проверяет доступ пользователя к записям
func (s *TaskListService) CheckAccess(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "CheckAccess")
	defer span.End()

	return s.taskListRepository.CheckAccess(ctx, userID, ids)
}
