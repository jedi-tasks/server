package service

import (
	"context"

	"github.com/google/uuid"

	"server/internal/modules/task/repository/taskDDL"
)

// HasLiveTasksInList проверяет, есть ли неудаленные задачи в списке
func (s *TaskService) HasLiveTasksInList(ctx context.Context, id uuid.UUID) (bool, error) {
	ctx, span := tracer.Start(ctx, "HasLiveTasksInList")
	defer span.End()

	return s.taskRepository.HasLiveTasks(ctx, taskDDL.ColumnListID, id)
}
