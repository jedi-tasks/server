package service

import (
	"context"

	"server/internal/modules/tag/model"
)

// GetTags возвращает записи пользователя по фильтрам
func (s *TagService) GetTags(ctx context.Context, req model.GetTagsReq) ([]model.Tag, error) {
	ctx, span := tracer.Start(ctx, "GetTags")
	defer span.End()

	return s.tagRepository.GetTags(ctx, req)
}
