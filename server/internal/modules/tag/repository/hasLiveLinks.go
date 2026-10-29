package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/tag/repository/taskToTagDDL"
)

// HasLiveLinks проверяет, есть ли неудаленные связи задач и тегов, ссылающиеся на запись.
// column — taskToTagDDL.ColumnTaskID или taskToTagDDL.ColumnTagID
func (r *TagRepository) HasLiveLinks(ctx context.Context, column string, referenceID uuid.UUID) (exists bool, err error) {
	ctx, span := tracer.Start(ctx, "HasLiveLinks")
	defer span.End()

	return exists, r.db.Get(ctx, &exists, sq.Select().Column(sq.Expr("EXISTS(?)", sq.
		Select("1").
		From(taskToTagDDL.Table).
		Where(sq.Eq{
			column:                       referenceID,
			taskToTagDDL.ColumnDeletedAt: nil,
		}),
	)))
}
