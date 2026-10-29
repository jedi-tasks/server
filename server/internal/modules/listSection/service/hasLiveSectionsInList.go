package service

import (
	"context"

	"github.com/google/uuid"
)

// HasLiveSectionsInList проверяет, есть ли в списке неудаленные разделы
func (s *ListSectionService) HasLiveSectionsInList(ctx context.Context, listID uuid.UUID) (bool, error) {
	ctx, span := tracer.Start(ctx, "HasLiveSectionsInList")
	defer span.End()

	return s.listSectionRepository.HasLiveSectionsInList(ctx, listID)
}
