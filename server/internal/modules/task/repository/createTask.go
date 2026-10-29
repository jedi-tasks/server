package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"server/internal/modules/task/repository/model"
	"server/internal/modules/task/repository/taskDDL"
)

// CreateTask создает запись. Повтор того же запроса из офлайн-очереди ничего не меняет:
// возвращается false, если запись с таким идентификатором уже есть
func (r *TaskRepository) CreateTask(ctx context.Context, req model.CreateTaskReq) (created bool, err error) {
	ctx, span := tracer.Start(ctx, "CreateTask")
	defer span.End()

	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Insert(taskDDL.Table).
		SetMap(map[string]any{
			taskDDL.ColumnID:              req.ID,
			taskDDL.ColumnUserID:          req.UserID,
			taskDDL.ColumnListID:          req.ListID,
			taskDDL.ColumnSectionID:       req.SectionID,
			taskDDL.ColumnTitle:           req.Title,
			taskDDL.ColumnNotes:           req.Notes,
			taskDDL.ColumnPriorityID:      req.PriorityID,
			taskDDL.ColumnPlannedDay:      req.PlannedDay,
			taskDDL.ColumnDeadlineDay:     req.DeadlineDay,
			taskDDL.ColumnDurationMinutes: req.DurationMinutes,
			taskDDL.ColumnManualOrder:     req.ManualOrder,
			taskDDL.ColumnExternalID:      req.ExternalID,
			taskDDL.ColumnCreatedAt:       req.CreatedAt,
			taskDDL.ColumnUpdatedAt:       req.CreatedAt,
		}).
		Suffix("ON CONFLICT ("+taskDDL.ColumnID+") DO NOTHING"),
	)
	if err != nil {
		return false, err
	}

	return rows != 0, nil
}
