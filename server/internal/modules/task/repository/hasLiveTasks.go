package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/task/repository/taskDDL"
)

// HasLiveTasks проверяет, есть ли неудаленные задачи, ссылающиеся на запись. column — taskDDL.ColumnListID,
// taskDDL.ColumnSectionID или taskDDL.ColumnPriorityID
func (r *TaskRepository) HasLiveTasks(ctx context.Context, column string, referenceID uuid.UUID) (exists bool, err error) {
	ctx, span := tracer.Start(ctx, "HasLiveTasks")
	defer span.End()

	return exists, r.db.Get(ctx, &exists, sq.Select().Column(sq.Expr("EXISTS(?)", sq.
		Select("1").
		From(taskDDL.Table).
		Where(sq.Eq{
			column:                  referenceID,
			taskDDL.ColumnDeletedAt: nil,
		}),
	)))
}
