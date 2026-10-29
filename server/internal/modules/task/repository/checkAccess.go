package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"pkg/slices"
	"server/internal/modules/task/repository/taskDDL"
	"server/internal/utils/errors"
)

// CheckAccess проверяет, что все записи существуют, не удалены и принадлежат пользователю
func (r *TaskRepository) CheckAccess(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "CheckAccess")
	defer span.End()

	ids = slices.GetUniqueFields(ids, func(id uuid.UUID) uuid.UUID { return id })

	var count int
	if err := r.db.Get(ctx, &count, sq.
		Select("COUNT(*)").
		From(taskDDL.Table).
		Where(sq.Eq{
			taskDDL.ColumnID:        ids,
			taskDDL.ColumnUserID:    userID,
			taskDDL.ColumnDeletedAt: nil,
		}),
	); err != nil {
		return err
	}

	if count != len(ids) {
		return errors.Forbidden.New("You don't have access to some of the requested tasks").
			WithContextParams(ctx).
			WithParams("UserID", userID, "IDs", ids)
	}

	return nil
}
