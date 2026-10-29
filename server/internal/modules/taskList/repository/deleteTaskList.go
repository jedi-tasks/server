package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/taskList/repository/taskListDDL"
)

// DeleteTaskList мягко удаляет запись, если удаление новее последнего изменения. Возвращает false,
// если запись уже удалена или на сервере более новая версия
func (r *TaskListRepository) DeleteTaskList(ctx context.Context, id uuid.UUID, updatedAt time.Time) (deleted bool, err error) {
	ctx, span := tracer.Start(ctx, "DeleteTaskList")
	defer span.End()

	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Update(taskListDDL.Table).
		SetMap(map[string]any{
			taskListDDL.ColumnDeletedAt: updatedAt,
			taskListDDL.ColumnUpdatedAt: updatedAt,
		}).
		Where(sq.Eq{taskListDDL.ColumnID: id, taskListDDL.ColumnDeletedAt: nil}).
		Where(sq.Lt{taskListDDL.ColumnUpdatedAt: updatedAt}),
	)
	if err != nil {
		return false, err
	}

	return rows != 0, nil
}
