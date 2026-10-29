package service

import (
	"context"

	"server/internal/modules/task/model"
)

// GetTasks возвращает записи пользователя по фильтрам
func (s *TaskService) GetTasks(ctx context.Context, req model.GetTasksReq) ([]model.Task, error) {
	ctx, span := tracer.Start(ctx, "GetTasks")
	defer span.End()

	return s.taskRepository.GetTasks(ctx, req)
}
