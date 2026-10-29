package service

import (
	"context"

	"server/internal/modules/tag/model"
)

// GetTasksToTags возвращает связи задач и тегов пользователя
func (s *TagService) GetTasksToTags(ctx context.Context, req model.GetTasksToTagsReq) ([]model.TaskToTag, error) {
	ctx, span := tracer.Start(ctx, "GetTasksToTags")
	defer span.End()

	return s.tagRepository.GetTasksToTags(ctx, req)
}
