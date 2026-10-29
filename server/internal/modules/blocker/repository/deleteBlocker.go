package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/blocker/repository/blockerDDL"
)

// DeleteBlocker мягко удаляет запись, если удаление новее последнего изменения. Возвращает false,
// если запись уже удалена или на сервере более новая версия
func (r *BlockerRepository) DeleteBlocker(ctx context.Context, id uuid.UUID, updatedAt time.Time) (deleted bool, err error) {
	ctx, span := tracer.Start(ctx, "DeleteBlocker")
	defer span.End()

	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Update(blockerDDL.Table).
		SetMap(map[string]any{
			blockerDDL.ColumnDeletedAt: updatedAt,
			blockerDDL.ColumnUpdatedAt: updatedAt,
		}).
		Where(sq.Eq{blockerDDL.ColumnID: id, blockerDDL.ColumnDeletedAt: nil}).
		Where(sq.Lt{blockerDDL.ColumnUpdatedAt: updatedAt}),
	)
	if err != nil {
		return false, err
	}

	return rows != 0, nil
}
