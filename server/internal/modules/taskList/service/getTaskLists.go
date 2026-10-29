package service

import (
	"context"

	"server/internal/modules/taskList/model"
)

// GetTaskLists возвращает записи пользователя по фильтрам
func (s *TaskListService) GetTaskLists(ctx context.Context, req model.GetTaskListsReq) ([]model.TaskList, error) {
	ctx, span := tracer.Start(ctx, "GetTaskLists")
	defer span.End()

	return s.taskListRepository.GetTaskLists(ctx, req)
}
