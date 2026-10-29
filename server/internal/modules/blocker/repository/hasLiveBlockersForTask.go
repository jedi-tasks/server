package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/blocker/repository/blockerDDL"
)

// HasLiveBlockersForTask проверяет, есть ли неудаленные блокеры, где задача заблокирована или блокирует
func (r *BlockerRepository) HasLiveBlockersForTask(ctx context.Context, taskID uuid.UUID) (exists bool, err error) {
	ctx, span := tracer.Start(ctx, "HasLiveBlockersForTask")
	defer span.End()

	return exists, r.db.Get(ctx, &exists, sq.Select().Column(sq.Expr("EXISTS(?)", sq.
		Select("1").
		From(blockerDDL.Table).
		Where(sq.Or{
			sq.Eq{blockerDDL.ColumnTaskID: taskID},
			sq.Eq{blockerDDL.ColumnBlockingTaskID: taskID},
		}).
		Where(sq.Eq{blockerDDL.ColumnDeletedAt: nil}),
	)))
}
