package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/task/repository/taskDDL"
)

// SetTaskCompleted выставляет или снимает время выполнения задачи, если изменение новее сохраненного.
// completedAt == nil — вернуть в работу. Возвращает false, если на сервере уже более новая версия
func (r *TaskRepository) SetTaskCompleted(ctx context.Context, id uuid.UUID, completedAt *time.Time, updatedAt time.Time) (updated bool, err error) {
	ctx, span := tracer.Start(ctx, "SetTaskCompleted")
	defer span.End()

	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Update(taskDDL.Table).
		SetMap(map[string]any{
			taskDDL.ColumnCompletedAt: completedAt,
			taskDDL.ColumnUpdatedAt:   updatedAt,
		}).
		Where(sq.Eq{taskDDL.ColumnID: id, taskDDL.ColumnDeletedAt: nil}).
		Where(sq.Lt{taskDDL.ColumnUpdatedAt: updatedAt}),
	)
	if err != nil {
		return false, err
	}

	return rows != 0, nil
}
