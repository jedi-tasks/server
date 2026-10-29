package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"pkg/slices"
	"server/internal/modules/priority/repository/priorityDDL"
	"server/internal/utils/errors"
)

// CheckAccess проверяет, что все записи существуют, не удалены и принадлежат пользователю
func (r *PriorityRepository) CheckAccess(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "CheckAccess")
	defer span.End()

	ids = slices.GetUniqueFields(ids, func(id uuid.UUID) uuid.UUID { return id })

	var count int
	if err := r.db.Get(ctx, &count, sq.
		Select("COUNT(*)").
		From(priorityDDL.Table).
		Where(sq.Eq{
			priorityDDL.ColumnID:        ids,
			priorityDDL.ColumnUserID:    userID,
			priorityDDL.ColumnDeletedAt: nil,
		}),
	); err != nil {
		return err
	}

	if count != len(ids) {
		return errors.Forbidden.New("You don't have access to some of the requested priorities").
			WithContextParams(ctx).
			WithParams("UserID", userID, "IDs", ids)
	}

	return nil
}
