package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"server/internal/modules/taskList/repository/model"
	"server/internal/modules/taskList/repository/taskListDDL"
)

// CreateTaskList создает запись. Повтор того же запроса из офлайн-очереди ничего не меняет:
// возвращается false, если запись с таким идентификатором уже есть
func (r *TaskListRepository) CreateTaskList(ctx context.Context, req model.CreateTaskListReq) (created bool, err error) {
	ctx, span := tracer.Start(ctx, "CreateTaskList")
	defer span.End()

	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Insert(taskListDDL.Table).
		SetMap(map[string]any{
			taskListDDL.ColumnID:        req.ID,
			taskListDDL.ColumnUserID:    req.UserID,
			taskListDDL.ColumnName:      req.Name,
			taskListDDL.ColumnColorHex:  req.ColorHex,
			taskListDDL.ColumnIcon:      req.Icon,
			taskListDDL.ColumnSortOrder: req.SortOrder,
			taskListDDL.ColumnSortMode:  req.SortMode,
			taskListDDL.ColumnCreatedAt: req.CreatedAt,
			taskListDDL.ColumnUpdatedAt: req.CreatedAt,
		}).
		Suffix("ON CONFLICT ("+taskListDDL.ColumnID+") DO NOTHING"),
	)
	if err != nil {
		return false, err
	}

	return rows != 0, nil
}
