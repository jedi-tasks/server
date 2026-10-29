package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/durationOption/repository/durationOptionDDL"
)

// DeleteDurationOption мягко удаляет запись, если удаление новее последнего изменения. Возвращает false,
// если запись уже удалена или на сервере более новая версия
func (r *DurationOptionRepository) DeleteDurationOption(ctx context.Context, id uuid.UUID, updatedAt time.Time) (deleted bool, err error) {
	ctx, span := tracer.Start(ctx, "DeleteDurationOption")
	defer span.End()

	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Update(durationOptionDDL.Table).
		SetMap(map[string]any{
			durationOptionDDL.ColumnDeletedAt: updatedAt,
			durationOptionDDL.ColumnUpdatedAt: updatedAt,
		}).
		Where(sq.Eq{durationOptionDDL.ColumnID: id, durationOptionDDL.ColumnDeletedAt: nil}).
		Where(sq.Lt{durationOptionDDL.ColumnUpdatedAt: updatedAt}),
	)
	if err != nil {
		return false, err
	}

	return rows != 0, nil
}
