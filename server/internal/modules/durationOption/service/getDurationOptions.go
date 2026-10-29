package service

import (
	"context"

	"server/internal/modules/durationOption/model"
)

// GetDurationOptions возвращает записи пользователя по фильтрам
func (s *DurationOptionService) GetDurationOptions(ctx context.Context, req model.GetDurationOptionsReq) ([]model.DurationOption, error) {
	ctx, span := tracer.Start(ctx, "GetDurationOptions")
	defer span.End()

	return s.durationOptionRepository.GetDurationOptions(ctx, req)
}
