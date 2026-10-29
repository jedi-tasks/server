package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/listSection/repository/listSectionDDL"
)

// HasLiveSectionsInList проверяет, есть ли в списке неудаленные разделы
func (r *ListSectionRepository) HasLiveSectionsInList(ctx context.Context, listID uuid.UUID) (exists bool, err error) {
	ctx, span := tracer.Start(ctx, "HasLiveSectionsInList")
	defer span.End()

	return exists, r.db.Get(ctx, &exists, sq.Select().Column(sq.Expr("EXISTS(?)", sq.
		Select("1").
		From(listSectionDDL.Table).
		Where(sq.Eq{
			listSectionDDL.ColumnListID:    listID,
			listSectionDDL.ColumnDeletedAt: nil,
		}),
	)))
}
