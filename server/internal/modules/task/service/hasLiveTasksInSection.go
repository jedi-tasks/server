package service

import (
	"context"

	"github.com/google/uuid"

	"server/internal/modules/task/repository/taskDDL"
)

// HasLiveTasksInSection проверяет, есть ли неудаленные задачи в разделе
func (s *TaskService) HasLiveTasksInSection(ctx context.Context, id uuid.UUID) (bool, error) {
	ctx, span := tracer.Start(ctx, "HasLiveTasksInSection")
	defer span.End()

	return s.taskRepository.HasLiveTasks(ctx, taskDDL.ColumnSectionID, id)
}
