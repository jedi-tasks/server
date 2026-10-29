package service

import (
	"context"

	"github.com/google/uuid"
)

// GetInboxID возвращает идентификатор системного списка Inbox пользователя
func (s *TaskListService) GetInboxID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	ctx, span := tracer.Start(ctx, "GetInboxID")
	defer span.End()

	return s.taskListRepository.GetInboxID(ctx, userID)
}
