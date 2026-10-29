package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/tag/repository/tagDDL"
)

// DeleteTag мягко удаляет запись, если удаление новее последнего изменения. Возвращает false,
// если запись уже удалена или на сервере более новая версия
func (r *TagRepository) DeleteTag(ctx context.Context, id uuid.UUID, updatedAt time.Time) (deleted bool, err error) {
	ctx, span := tracer.Start(ctx, "DeleteTag")
	defer span.End()

	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Update(tagDDL.Table).
		SetMap(map[string]any{
			tagDDL.ColumnDeletedAt: updatedAt,
			tagDDL.ColumnUpdatedAt: updatedAt,
		}).
		Where(sq.Eq{tagDDL.ColumnID: id, tagDDL.ColumnDeletedAt: nil}).
		Where(sq.Lt{tagDDL.ColumnUpdatedAt: updatedAt}),
	)
	if err != nil {
		return false, err
	}

	return rows != 0, nil
}
