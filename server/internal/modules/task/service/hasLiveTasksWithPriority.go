package service

import (
	"context"

	"github.com/google/uuid"

	"server/internal/modules/task/repository/taskDDL"
)

// HasLiveTasksWithPriority проверяет, есть ли неудаленные задачи с приоритетом
func (s *TaskService) HasLiveTasksWithPriority(ctx context.Context, id uuid.UUID) (bool, error) {
	ctx, span := tracer.Start(ctx, "HasLiveTasksWithPriority")
	defer span.End()

	return s.taskRepository.HasLiveTasks(ctx, taskDDL.ColumnPriorityID, id)
}
