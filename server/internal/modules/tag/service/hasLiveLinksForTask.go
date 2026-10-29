package service

import (
	"context"

	"github.com/google/uuid"

	"server/internal/modules/tag/repository/taskToTagDDL"
)

// HasLiveLinksForTask проверяет, привязаны ли к задаче теги неудаленными связями
func (s *TagService) HasLiveLinksForTask(ctx context.Context, taskID uuid.UUID) (bool, error) {
	ctx, span := tracer.Start(ctx, "HasLiveLinksForTask")
	defer span.End()

	return s.tagRepository.HasLiveLinks(ctx, taskToTagDDL.ColumnTaskID, taskID)
}
