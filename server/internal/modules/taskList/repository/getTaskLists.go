package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"pkg/ddlHelper"
	"server/internal/modules/taskList/model"
	"server/internal/modules/taskList/repository/taskListDDL"
)

// GetTaskLists возвращает неудаленные записи пользователя по фильтрам
func (r *TaskListRepository) GetTaskLists(ctx context.Context, req model.GetTaskListsReq) (res []model.TaskList, err error) {
	ctx, span := tracer.Start(ctx, "GetTaskLists")
	defer span.End()

	q := sq.
		Select(ddlHelper.SelectAll).
		From(taskListDDL.Table).
		Where(sq.Eq{
			taskListDDL.ColumnUserID:    req.Necessary.UserID,
			taskListDDL.ColumnDeletedAt: nil,
		})

	// Добавляем фильтры
	if len(req.IDs) != 0 {
		q = q.Where(sq.Eq{taskListDDL.ColumnID: req.IDs})
	}

	return res, r.db.Select(ctx, &res, q)
}
