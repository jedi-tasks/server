package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"server/internal/modules/task/repository/taskDDL"
)

// ClearReference сбрасывает ссылку на удаленный раздел или приоритет у задач пользователя и возвращает
// идентификаторы измененных задач. column — taskDDL.ColumnSectionID или taskDDL.ColumnPriorityID
func (r *TaskRepository) ClearReference(ctx context.Context, userID uuid.UUID, column string, refID uuid.UUID, updatedAt time.Time) (ids []uuid.UUID, err error) {
	ctx, span := tracer.Start(ctx, "ClearReference")
	defer span.End()

	return ids, r.db.Select(ctx, &ids, sq.
		Update(taskDDL.Table).
		SetMap(map[string]any{
			column:                  nil,
			taskDDL.ColumnUpdatedAt: sq.Expr("GREATEST("+taskDDL.ColumnUpdatedAt+", ?)", updatedAt),
		}).
		Where(sq.Eq{
			taskDDL.ColumnUserID:    userID,
			column:                  refID,
			taskDDL.ColumnDeletedAt: nil,
		}).
		Suffix("RETURNING "+taskDDL.ColumnID),
	)
}
