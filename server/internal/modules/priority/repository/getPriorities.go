package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/ddlHelper"
	"server/internal/modules/priority/model"
	"server/internal/modules/priority/repository/priorityDDL"
)

// GetPriorities возвращает неудаленные записи пользователя по фильтрам
func (r *PriorityRepository) GetPriorities(ctx context.Context, req model.GetPrioritiesReq) (res []model.Priority, err error) {
	ctx, span := tracer.Start(ctx, "GetPriorities")
	defer span.End()

	q := sq.
		Select(ddlHelper.SelectAll).
		From(priorityDDL.Table).
		Where(sq.Eq{
			priorityDDL.ColumnUserID:    req.Necessary.UserID,
			priorityDDL.ColumnDeletedAt: nil,
		})

	// Добавляем фильтры
	if len(req.IDs) != 0 {
		q = q.Where(sq.Eq{priorityDDL.ColumnID: req.IDs})
	}

	return res, r.db.Select(ctx, &res, q)
}
