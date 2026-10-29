package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"server/internal/modules/taskList/repository/model"
	"server/internal/modules/taskList/repository/taskListDDL"
)

// CreateInbox создает системный список Inbox
func (r *TaskListRepository) CreateInbox(ctx context.Context, req model.CreateTaskListReq) error {
	ctx, span := tracer.Start(ctx, "CreateInbox")
	defer span.End()

	return r.db.Exec(ctx, sq.
		Insert(taskListDDL.Table).
		SetMap(map[string]any{
			taskListDDL.ColumnID:        req.ID,
			taskListDDL.ColumnUserID:    req.UserID,
			taskListDDL.ColumnName:      req.Name,
			taskListDDL.ColumnIsInbox:   true,
			taskListDDL.ColumnColorHex:  req.ColorHex,
			taskListDDL.ColumnIcon:      req.Icon,
			taskListDDL.ColumnSortOrder: req.SortOrder,
			taskListDDL.ColumnSortMode:  req.SortMode,
			taskListDDL.ColumnCreatedAt: req.CreatedAt,
			taskListDDL.ColumnUpdatedAt: req.CreatedAt,
		}),
	)
}
