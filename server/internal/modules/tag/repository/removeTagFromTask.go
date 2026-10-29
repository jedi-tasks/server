package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/tag/repository/taskToTagDDL"
)

// RemoveTagFromTask мягко удаляет связь задачи и тега, если отвязка новее последнего изменения связи.
// Возвращает false, если связи нет или на сервере более новая версия
func (r *TagRepository) RemoveTagFromTask(ctx context.Context, taskID, tagID uuid.UUID, updatedAt time.Time) (removed bool, err error) {
	ctx, span := tracer.Start(ctx, "RemoveTagFromTask")
	defer span.End()

	rows, err := r.db.ExecWithRowsAffected(ctx, sq.
		Update(taskToTagDDL.Table).
		SetMap(map[string]any{
			taskToTagDDL.ColumnDeletedAt: updatedAt,
			taskToTagDDL.ColumnUpdatedAt: updatedAt,
		}).
		Where(sq.Eq{
			taskToTagDDL.ColumnTaskID:    taskID,
			taskToTagDDL.ColumnTagID:     tagID,
			taskToTagDDL.ColumnDeletedAt: nil,
		}).
		Where(sq.Lt{taskToTagDDL.ColumnUpdatedAt: updatedAt}),
	)
	if err != nil {
		return false, err
	}

	return rows != 0, nil
}
