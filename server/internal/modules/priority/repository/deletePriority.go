package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/priority/repository/priorityDDL"
)

// DeletePriority мягко удаляет запись, если удаление новее последнего изменения. Возвращает false,
// если запись уже удалена или на сервере более новая версия
func (r *PriorityRepository) DeletePriority(ctx context.Context, id uuid.UUID, updatedAt time.Time) (deleted bool, err error) {
	ctx, span := tracer.Start(ctx, "DeletePriority")
	defer span.End()

	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Update(priorityDDL.Table).
		SetMap(map[string]any{
			priorityDDL.ColumnDeletedAt: updatedAt,
			priorityDDL.ColumnUpdatedAt: updatedAt,
		}).
		Where(sq.Eq{priorityDDL.ColumnID: id, priorityDDL.ColumnDeletedAt: nil}).
		Where(sq.Lt{priorityDDL.ColumnUpdatedAt: updatedAt}),
	)
	if err != nil {
		return false, err
	}

	return rows != 0, nil
}
