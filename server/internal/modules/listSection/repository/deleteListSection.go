package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/listSection/repository/listSectionDDL"
)

// DeleteListSection мягко удаляет запись, если удаление новее последнего изменения. Возвращает false,
// если запись уже удалена или на сервере более новая версия
func (r *ListSectionRepository) DeleteListSection(ctx context.Context, id uuid.UUID, updatedAt time.Time) (deleted bool, err error) {
	ctx, span := tracer.Start(ctx, "DeleteListSection")
	defer span.End()

	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Update(listSectionDDL.Table).
		SetMap(map[string]any{
			listSectionDDL.ColumnDeletedAt: updatedAt,
			listSectionDDL.ColumnUpdatedAt: updatedAt,
		}).
		Where(sq.Eq{listSectionDDL.ColumnID: id, listSectionDDL.ColumnDeletedAt: nil}).
		Where(sq.Lt{listSectionDDL.ColumnUpdatedAt: updatedAt}),
	)
	if err != nil {
		return false, err
	}

	return rows != 0, nil
}
