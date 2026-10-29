package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/ddlHelper"
	"server/internal/modules/blocker/model"
	"server/internal/modules/blocker/repository/blockerDDL"
)

// GetBlockers возвращает неудаленные записи пользователя по фильтрам
func (r *BlockerRepository) GetBlockers(ctx context.Context, req model.GetBlockersReq) (res []model.Blocker, err error) {
	ctx, span := tracer.Start(ctx, "GetBlockers")
	defer span.End()

	q := sq.
		Select(ddlHelper.SelectAll).
		From(blockerDDL.Table).
		Where(sq.Eq{
			blockerDDL.ColumnUserID:    req.Necessary.UserID,
			blockerDDL.ColumnDeletedAt: nil,
		})

	// Добавляем фильтры
	if len(req.IDs) != 0 {
		q = q.Where(sq.Eq{blockerDDL.ColumnID: req.IDs})
	}
	if len(req.TaskIDs) != 0 {
		q = q.Where(sq.Eq{blockerDDL.ColumnTaskID: req.TaskIDs})
	}

	if len(req.BlockingTaskIDs) != 0 {
		q = q.Where(sq.Eq{blockerDDL.ColumnBlockingTaskID: req.BlockingTaskIDs})
	}
	if req.OnlyActive {
		q = q.Where(sq.Eq{blockerDDL.ColumnResolvedAt: nil})
	}

	return res, r.db.Select(ctx, &res, q)
}
