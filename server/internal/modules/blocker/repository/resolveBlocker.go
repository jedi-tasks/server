package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/blocker/repository/blockerDDL"
)

// ResolveBlocker снимает блокер, если он еще жив и снятие новее последнего изменения. Возвращает false,
// если блокер уже снят или на сервере более новая версия
func (r *BlockerRepository) ResolveBlocker(ctx context.Context, id uuid.UUID, resolvedAt time.Time) (resolved bool, err error) {
	ctx, span := tracer.Start(ctx, "ResolveBlocker")
	defer span.End()

	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Update(blockerDDL.Table).
		SetMap(map[string]any{
			blockerDDL.ColumnResolvedAt: resolvedAt,
			blockerDDL.ColumnUpdatedAt:  sq.Expr("GREATEST("+blockerDDL.ColumnUpdatedAt+", ?)", resolvedAt),
		}).
		Where(sq.Eq{blockerDDL.ColumnID: id, blockerDDL.ColumnResolvedAt: nil, blockerDDL.ColumnDeletedAt: nil}),
	)
	if err != nil {
		return false, err
	}

	return rows != 0, nil
}
