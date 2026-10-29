package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/ddlHelper"
	"server/internal/modules/task/model"
	"server/internal/modules/task/repository/taskDDL"
)

// GetTasks возвращает неудаленные записи пользователя по фильтрам
func (r *TaskRepository) GetTasks(ctx context.Context, req model.GetTasksReq) (res []model.Task, err error) {
	ctx, span := tracer.Start(ctx, "GetTasks")
	defer span.End()

	q := sq.
		Select(ddlHelper.SelectAll).
		From(taskDDL.Table).
		Where(sq.Eq{
			taskDDL.ColumnUserID:    req.Necessary.UserID,
			taskDDL.ColumnDeletedAt: nil,
		})

	// Добавляем фильтры
	if len(req.IDs) != 0 {
		q = q.Where(sq.Eq{taskDDL.ColumnID: req.IDs})
	}
	if len(req.ListIDs) != 0 {
		q = q.Where(sq.Eq{taskDDL.ColumnListID: req.ListIDs})
	}
	if req.Completed != nil {
		if *req.Completed {
			q = q.Where(sq.NotEq{taskDDL.ColumnCompletedAt: nil})
		} else {
			q = q.Where(sq.Eq{taskDDL.ColumnCompletedAt: nil})
		}
	}

	return res, r.db.Select(ctx, &res, q)
}
