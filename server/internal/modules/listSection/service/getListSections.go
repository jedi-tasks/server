package service

import (
	"context"

	"server/internal/modules/listSection/model"
)

// GetListSections возвращает записи пользователя по фильтрам
func (s *ListSectionService) GetListSections(ctx context.Context, req model.GetListSectionsReq) ([]model.ListSection, error) {
	ctx, span := tracer.Start(ctx, "GetListSections")
	defer span.End()

	return s.listSectionRepository.GetListSections(ctx, req)
}
